## 一、结构体和变量
```markdown
struct Pool {
    address stTokenAddress; // 质押代币的地址
    uint256 poolWeight; // 不同资金池所占的权重
                        // 示例：ETH 池权重 100，USDT 池权重 50 → ETH 池获得 2/3 的奖励
    uint256 lastRewardBlock; // 最后一次奖励的区块编号
    uint256 accMetaNodePerST; // 质押 1个ETH经过1个区块高度，能拿到几个MetaNode奖励币
    uint256 stTokenAmount; // 质押的代币数量
    uint256 minDepositAmount; // 最小质押代币数量
    uint256 unstakeLockedBlocks; // 解质押锁定的区块高度，即经过几个区块后才可以提取质押代币
}
```

```markdown
struct UnstakeRequest {
    uint256 amount; // 用户取消质押的代币数量
    uint256 unlockBlocks; // 解质押的区块高度
}
```

```markdown
struct User { // 记录用户相对每个资金池 的质押记录
    uint256 stAmount; // 用户在当前资金池，质押的代币数量
    uint256 finishedMetaNode; // 用户在当前资金池，已经领取的 MetaNode 奖励数量
    uint256 pendingMetaNode; // 用户在当前资金池，当前可领取的 MetaNode 奖励数量
    UnstakeRequest[] requests; // 用户在当前资金池，取消质押的记录
}
```

```markdown
    uint256 public startBlock; // 质押开始区块高度
    uint256 public endBlock; // 质押结束区块高度
    uint256 public MetaNodePerBlock; // 每个区块高度，MetaNode 的奖励数量
    bool public withdrawPaused; // 是否暂停提现质押代币
    bool public claimPaused; // 是否暂停领取奖励
    IERC20 public MetaNode; // MetaNode 奖励代币地址
    uint256 public totalPoolWeight; // 所有资金池的权重总和
    Pool[] public pool; // 资金池列表
    // 资金池 id => 用户地址 => 用户信息
    mapping (uint256 => mapping (address => User)) public user; 
```

## 二、函数
```markdown
    function initialize(
        IERC20 _MetaNode,
        uint256 _startBlock,
        uint256 _endBlock,
        uint256 _MetaNodePerBlock
    ) public initializer { // 替代构造函数，用于可升级合约的初始化。
        require(_startBlock <= _endBlock && _MetaNodePerBlock > 0, "invalid parameters");

        __Pausable_init(); // 提供暂停某些功能的能力（如提现、领取奖励）
        __AccessControl_init(); // 基于角色的权限控制
        __UUPSUpgradeable_init(); // 实现 UUPS 代理升级模式。
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(UPGRADE_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);

        setMetaNode(_MetaNode); //设置奖励代币

        startBlock = _startBlock; // 质押开始区块高度
        endBlock = _endBlock; // 质押结束区块高度
        MetaNodePerBlock = _MetaNodePerBlock; // 每个区块高度，MetaNode 的奖励数量

    }
```

```markdown
    function setMetaNode(IERC20 _MetaNode) public onlyRole(ADMIN_ROLE) {
        MetaNode = _MetaNode;
        emit SetMetaNode(MetaNode);
    }
```

```markdown
    //只有拥有 UPGRADE_ROLE 的地址才能升级合约实现。
    function _authorizeUpgrade(address newImplementation)
        internal
        onlyRole(UPGRADE_ROLE)
        override
    { }
```

### 核心函数
#### 1、创建质押池
+ **添加资金池**：
    - 管理员可以配置多个质押池，每个池有不同的质押代币和权重。
    - 第一个池必须是 ETH（`address(0)`），对应 `ETH_PID = 0`
    - 所有质押池的添加，只能由admin权限的管理员进行。

```markdown
function addPool(address _stTokenAddress, uint256 _poolWeight, uint256 _minDepositAmount, 
                 uint256 _unstakeLockedBlocks, bool _withUpdate) 
  public onlyRole(ADMIN_ROLE) {
  
    // 第一个池必须是 ETH（address(0)），对应 ETH_PID = 0
    if (pool.length > 0) {
        require(_stTokenAddress != address(0x0), "invalid staking token address");
    } else {
        require(_stTokenAddress == address(0x0), "invalid staking token address");
    }

    //质押锁区块必须大于0
    require(_unstakeLockedBlocks > 0, "invalid withdraw locked blocks");
    //质押活动未结束才允许添加池
    require(block.number < endBlock, "Already ended");

    if (_withUpdate) {
        //批量更新所有资金池状态
        //目的是保证在加入新池前，旧池的 accMetaNodePerST 已经正确计算到当前区块，避免奖励分配偏差。
        massUpdatePools();
    }

    //设置最后一次奖励的区块编号
    uint256 lastRewardBlock = block.number > startBlock ? block.number : startBlock;
    //计算所有质押池的权重
    totalPoolWeight = totalPoolWeight + _poolWeight;

    pool.push(Pool({
        stTokenAddress: _stTokenAddress,
        poolWeight: _poolWeight,
        lastRewardBlock: lastRewardBlock,
        accMetaNodePerST: 0,
        stTokenAmount: 0,
        minDepositAmount: _minDepositAmount,
        unstakeLockedBlocks: _unstakeLockedBlocks
    }));

    emit AddPool(_stTokenAddress, _poolWeight, lastRewardBlock, _minDepositAmount, _unstakeLockedBlocks);
}
```

```markdown
function massUpdatePools() public {
    uint256 length = pool.length;
    for (uint256 pid = 0; pid < length; pid++) {
        //循环更新每个质押池
        updatePool(pid);
    }
}
```

+ **更新指定质押池的奖励状态：**
    - 计算从上一次奖励更新到现在（或到结束块）之间，该池应得的 MetaNode 奖励。
    - 更新`accMetaNodePerST`：每单位质押代币累积可得的 MetaNode 数量
    - 更新`lastRewardBlock`：更新为当前区块

```markdown
function updatePool(uint256 _pid) public checkPid(_pid) {
    Pool storage pool_ = pool[_pid];

    //如果还没到奖励时间，直接返回
    if (block.number <= pool_.lastRewardBlock) {
        return;
    }

    //计算奖励发放周期内应得奖励总量, 奖励总量 = (有效区块数) × MetaNodePerBlock
    //totalMetaNode = 奖励总量*质押池权重
    (bool success1, uint256 totalMetaNode) 
    = getMultiplier(pool_.lastRewardBlock, block.number).tryMul(pool_.poolWeight);
    require(success1, "overflow");
    
    //再除以所有资源池的权重和，将总奖励按总权重比例缩放
    //totalMetaNode 现在就表示：这个池在这段时间内应得的 MetaNode 总量
    (success1, totalMetaNode) = totalMetaNode.tryDiv(totalPoolWeight);
    require(success1, "overflow");

    //如果池中有质押代币，计算并更新 accMetaNodePerST-每单位质押代币累积可得的 MetaNode 数量
    uint256 stSupply = pool_.stTokenAmount;
    if (stSupply > 0) {
        //将奖励数量放大 1e18 倍，用于高精度计算（避免整数除法丢失精度）
        (bool success2, uint256 totalMetaNode_) = totalMetaNode.tryMul(1 ether);
        require(success2, "overflow");

        //获得：这段区间内，每单位质押代币应得的 MetaNode 数量（带 18 位小数）
        (success2, totalMetaNode_) = totalMetaNode_.tryDiv(stSupply);
        require(success2, "overflow");

        //累加到累积奖励因子中
        (bool success3, uint256 accMetaNodePerST) 
        = pool_.accMetaNodePerST.tryAdd(totalMetaNode_);
        require(success3, "overflow");
        
        pool_.accMetaNodePerST = accMetaNodePerST;
    }

    pool_.lastRewardBlock = block.number;

    emit UpdatePool(_pid, pool_.lastRewardBlock, totalMetaNode);
}
```

+ **计算奖励发放周期内应得奖励总量**：
    - 奖励总量 = (有效区块数) × MetaNodePerBlock

```markdown
function getMultiplier(uint256 _from, uint256 _to) public view returns(uint256 multiplier) {
    require(_from <= _to, "invalid block");
    
    if (_from < startBlock) {_from = startBlock;}
    if (_to > endBlock) {_to = endBlock;}
    require(_from <= _to, "end block must be greater than start block");
    
    bool success;
    (success, multiplier) = (_to - _from).tryMul(MetaNodePerBlock);
    require(success, "multiplier overflow");
}
```

#### 2、质押
+ 用户质押ETH

```markdown
function depositETH() public whenNotPaused() payable {
    Pool storage pool_ = pool[ETH_PID];
    require(pool_.stTokenAddress == address(0x0), "invalid staking token address");

    uint256 _amount = msg.value;
    require(_amount >= pool_.minDepositAmount, "deposit amount is too small");

    _deposit(ETH_PID, _amount);
}
```

+ 用户质押ERC20代币
+ 注意：
    - 管理员要先创建对应代币的质押池
    - **用户需先调用 approve(...) 给本合约授权代币**

```markdown
function deposit(uint256 _pid, uint256 _amount) public whenNotPaused() checkPid(_pid) {
    require(_pid != 0, "deposit not support ETH staking");
    Pool storage pool_ = pool[_pid];
    require(_amount > pool_.minDepositAmount, "deposit amount is too small");

    if(_amount > 0) {
        IERC20(pool_.stTokenAddress).safeTransferFrom(msg.sender, address(this), _amount);
    }

    _deposit(_pid, _amount);
}
```

+ **最核心函数**，处理质押金额，完成的功能：**奖励结算器 + 状态同步器 + 权益更新器。**
+ 用户完成本次质押后，系统会认为：**从现在起，该用户的新质押量将参与未来奖励分配，而过去的历史奖励已结算完毕。**

```markdown
function _deposit(uint256 _pid, uint256 _amount) internal {
    //获取池和用户的数据引用,后续修改会永久保存,避免复制数据，节省 gas
    Pool storage pool_ = pool[_pid];
    User storage user_ = user[_pid][msg.sender];

    //更新指定质押池的奖励状态,确保 pool_.accMetaNodePerST 已经更新到当前区块。
    updatePool(_pid);

    //若用户已有质押
    if (user_.stAmount > 0) {
        //应得奖励 = 用户质押量 × 每单位代币累计奖励 - 已领取部分
        //        = stAmount × accMetaNodePerST / 1e18 - finishedMetaNode
        (bool success1, uint256 accST) = user_.stAmount.tryMul(pool_.accMetaNodePerST);
        require(success1, "user stAmount mul accMetaNodePerST overflow");
        
        (success1, accST) = accST.tryDiv(1 ether);
        require(success1, "accST div 1 ether overflow");
        
        (bool success2, uint256 pendingMetaNode_) = accST.trySub(user_.finishedMetaNode);
        require(success2, "accST sub finishedMetaNode overflow");

        //如果原来就有应得奖励，就加上本次的应得奖励
        if(pendingMetaNode_ > 0) {
            (bool success3, uint256 _pendingMetaNode) = user_.pendingMetaNode.tryAdd(pendingMetaNode_);
            require(success3, "user pendingMetaNode overflow");
            user_.pendingMetaNode = _pendingMetaNode;
        }
    }

    //增加用户质押币的数量
    if(_amount > 0) {
        (bool success4, uint256 stAmount) = user_.stAmount.tryAdd(_amount);
        require(success4, "user stAmount overflow");
        user_.stAmount = stAmount;
    }

    //增加质押池的总质押量
    (bool success5, uint256 stTokenAmount) = pool_.stTokenAmount.tryAdd(_amount);
    require(success5, "pool stTokenAmount overflow");
    pool_.stTokenAmount = stTokenAmount;

    // 重置用户的“已完成奖励”记录
    (bool success6, uint256 finishedMetaNode) = user_.stAmount.tryMul(pool_.accMetaNodePerST);
    require(success6, "user stAmount mul accMetaNodePerST overflow");

    (success6, finishedMetaNode) = finishedMetaNode.tryDiv(1 ether);
    require(success6, "finishedMetaNode div 1 ether overflow");

    user_.finishedMetaNode = finishedMetaNode;

    emit Deposit(msg.sender, _pid, _amount);
}
```

#### 3、申请解质押本金
+ 处理用户申请解质押（取款）：
    - 不是立即返还本金
    - 是创建一个**带锁定期的解质押请求**，用户需等待指定区块数后才能最终提取资金。

```markdown
function unstake(uint256 _pid, uint256 _amount) public whenNotPaused() checkPid(_pid) whenNotWithdrawPaused() {
    Pool storage pool_ = pool[_pid];
    User storage user_ = user[_pid][msg.sender];

    require(user_.stAmount >= _amount, "Not enough staking token balance");

    //更新指定质押池的奖励状态,确保 pool_.accMetaNodePerST 已经更新到当前区块。
    updatePool(_pid);

    //计算待领取奖励,即使用户只取本金，系统也会自动“结息”。
    uint256 pendingMetaNode_ 
    = user_.stAmount * pool_.accMetaNodePerST / (1 ether) - user_.finishedMetaNode;

    //累积待领取奖励
    if(pendingMetaNode_ > 0) {
        user_.pendingMetaNode = user_.pendingMetaNode + pendingMetaNode_;
    }

    if(_amount > 0) {
        //更新用户状态,立即减少用户的质押量（不再参与未来奖励分配）,但本金不会立刻返还
        user_.stAmount = user_.stAmount - _amount;
        //创建解质押请求
        user_.requests.push(UnstakeRequest({
            amount: _amount,
            unlockBlocks: block.number + pool_.unstakeLockedBlocks
        }));
    }

    //更新池的总质押量
    pool_.stTokenAmount = pool_.stTokenAmount - _amount;
    //重置用户的已完成奖励记录
    user_.finishedMetaNode = user_.stAmount * pool_.accMetaNodePerST / (1 ether);

    emit RequestUnstake(msg.sender, _pid, _amount);
}
```

```markdown
用户调用 unstake(_pid, _amount)
         │
         ▼
   检查余额 + 更新池状态
         │
         ▼
结算用户此前应得奖励 → pendingMetaNode
         │
         ▼
减少用户质押量（stAmount）
         │
         ▼
创建解质押请求（带解锁区块）
         │
         ▼
减少池总质押量
         │
         ▼
更新用户 finishedMetaNode
         │
         ▼
触发 RequestUnstake 事件
```

#### 4、提取解质押本金
+ 用户最终提取已解锁本金：
    - 用户可以批量领取所有已满足锁定期的解质押请求，并将代币安全地返还给用户。
    - 此函数不指定提取数量，而是自动提取所有已解锁的解质押请求。

```markdown
function withdraw(uint256 _pid) public whenNotPaused() checkPid(_pid) whenNotWithdrawPaused() {
    Pool storage pool_ = pool[_pid];
    User storage user_ = user[_pid][msg.sender];

    uint256 pendingWithdraw_;//已解锁的代币数量
    uint256 popNum_; //记录有多少个请求可以被移除
    //遍历解质押请求队列，找出所有已解锁的请求
    //user_.requests 是一个 UnstakeRequest[] 数组，按时间顺序排列（先进先出）
    for (uint256 i = 0; i < user_.requests.length; i++) {
        //未达到解锁条件
        if (user_.requests[i].unlockBlocks > block.number) {
            break;
        }
        pendingWithdraw_ = pendingWithdraw_ + user_.requests[i].amount;
        popNum_++;
    }

    //移动数组元素（高效删除前 N 个元素）
    //假设原数组有 5 个请求，其中前 2 个已解锁（popNum_ = 2）,将后面的元素整体前移 popNum_ 位。
    for (uint256 i = 0; i < user_.requests.length - popNum_; i++) {
        user_.requests[i] = user_.requests[i + popNum_];
    }

    //删除最后 popNum_ 个元素
    //pop() 从动态数组末尾移除一个元素。
    //执行 popNum_ 次后，所有已处理的请求都被清除。
    //这种“先移位再 pop”是 Solidity 中高效删除数组前缀元素的标准做法，比逐个删除更省 gas。
    for (uint256 i = 0; i < popNum_; i++) {
        user_.requests.pop();
    }

    //如果有待提取金额，执行转账
    if (pendingWithdraw_ > 0) {
        if (pool_.stTokenAddress == address(0x0)) {
            _safeETHTransfer(msg.sender, pendingWithdraw_);
        } else {
            IERC20(pool_.stTokenAddress).safeTransfer(msg.sender, pendingWithdraw_);
        }
    }

    emit Withdraw(msg.sender, _pid, pendingWithdraw_, block.number);
}
```

```markdown
用户调用 withdraw(_pid)
         │
         ▼
遍历 requests 数组
         │
         ▼
统计所有 unlockBlocks <= 当前区块的请求
         │
         ▼
累加 amount → pendingWithdraw_
         │
         ▼
将未解锁请求前移（覆盖已解锁的）
         │
         ▼
从末尾 pop() 掉已处理的请求
         │
         ▼
向用户转账 pendingWithdraw_
         │
         ▼
触发 Withdraw 事件
```

#### 5、领取累积奖励
+ 用户领取累积奖励：
    - 用户可以提取他们通过质押所应得的代币奖励。
    - 此函数不指定领取数量，而是自动领取**所有待发奖励**。

```markdown
function claim(uint256 _pid) public whenNotPaused() checkPid(_pid) whenNotClaimPaused() {
    Pool storage pool_ = pool[_pid];
    User storage user_ = user[_pid][msg.sender];

    //更新指定质押池的奖励状态,确保 pool_.accMetaNodePerST 已经更新到当前区块。
    updatePool(_pid);

    //计算总待领取奖励
    //总待领取奖励 = 截至目前理论上总共应得的 MetaNode 数量 
                    - 已计入账的部分 
                    + 之前未领取、已累积的奖励
    uint256 pendingMetaNode_ 
    = user_.stAmount * pool_.accMetaNodePerST / (1 ether) 
      - user_.finishedMetaNode 
      + user_.pendingMetaNode;

    //如果有待领取奖励，则发放并清零
    if(pendingMetaNode_ > 0) {
        user_.pendingMetaNode = 0;
        _safeMetaNodeTransfer(msg.sender, pendingMetaNode_);
    }

    //重置用户的已完成奖励记录
    user_.finishedMetaNode = user_.stAmount * pool_.accMetaNodePerST / (1 ether);

    emit Claim(msg.sender, _pid, pendingMetaNode_);
}
```

```markdown
用户调用 claim(_pid)
         │
         ▼
   更新池奖励状态（updatePool）
         │
         ▼
计算总奖励 = 当前理论总额 - 已完成额 + 待发额
         │
         ▼
如果 > 0，则转账给用户
         │
         ▼
清空 pendingMetaNode
         │
         ▼
重置 finishedMetaNode 为当前值
         │
         ▼
触发 Claim 事件
```

#### 6、安全转账
+ ETH安全转账

```markdown
function _safeETHTransfer(address _to, uint256 _amount) internal {
    (bool success, bytes memory data) = address(_to).call{
        value: _amount
    }("");

    require(success, "ETH transfer call failed");
    if (data.length > 0) {
        require(
            abi.decode(data, (bool)),
            "ETH transfer operation did not succeed"
        );
    }
}
```

+ ERC20代币安全转账

```markdown
function _safeMetaNodeTransfer(address _to, uint256 _amount) internal {
    uint256 MetaNodeBal = MetaNode.balanceOf(address(this));

    if (_amount > MetaNodeBal) {
        MetaNode.transfer(_to, MetaNodeBal);
    } else {
        MetaNode.transfer(_to, _amount);
    }
}
```

## 三、业务时序图
### 1、合约初始化
![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761738827317-32711e83-3579-4274-a3bf-9892cf56d299.png)

### 2、管理员更新质押池权重
![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761738879389-f01ac110-3420-48ed-adf2-c690a4adaec4.png)

### 3、用户质押 ETH
![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761738963795-8039e1b8-b083-42f1-93a6-c692ccd00cce.png)

### 4、用户质押 ERC20 代币
![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761738988983-0b7bc43a-8bd4-4f3a-8a46-eb0d2b0e0a27.png)

### 5、用户解质押与提现
![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761739017327-04a837d4-cb76-48c0-9725-78521319cbd7.png)

### 6、用户领取奖励
### ![](https://cdn.nlark.com/yuque/0/2025/png/40797156/1761739055120-b57d96be-0095-46a6-9d86-ce1d9b796467.png)


