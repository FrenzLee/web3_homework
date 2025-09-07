// SPDX-License-Identifier: MIT
pragma solidity ^0.8;
/*
作业 1：ERC20 代币
任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求： 
1. 合约包含以下标准 ERC20 功能：
2. balanceOf：查询账户余额。
3. transfer：转账。
4. approve 和 transferFrom：授权和代扣转账。
5. 使用 event 记录转账和授权操作。
6. 提供 mint 函数，允许合约所有者增发代币。
提示： 
● 使用 mapping 存储账户余额和授权信息。
● 使用 event 定义 Transfer 和 Approval 事件。
● 部署到sepolia 测试网，导入到自己的钱包。
*/

contract MyERC20 {

    //账户余额
    mapping(address => uint256) private _balances;
    //授权信息
    mapping(address => mapping(address => uint256)) private _allowances;

    string private _name;//币名
    string private _symbol;//币标志
    uint256 private _totalSupply;//币供应总数
    address private _owner;//发行者地址

    //定义事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);

    //构造函数
    constructor(string memory name_, string memory symbol_) {
        _name = name_;
        _symbol = symbol_;
        _owner = msg.sender;
    }

    //校验操作发起者只能是合约拥有者
    modifier onlyOwner() {
        require(_owner == msg.sender, "caller is not the owner");
        _;
    }

    //获取币名
    function name() public view virtual returns (string memory) {
        return _name;
    }

    //获取币标志
    function symbol() public view virtual returns (string memory) {
        return _symbol;
    }

    //获取精度
    function decimals() public view virtual returns (uint8) {
        return 18;
    }

    //获取币供应总数
    function totalSupply() public view virtual returns (uint256) {
        return _totalSupply;
    }


    //造币
    function mint(uint256 amount) public virtual onlyOwner {
        require(amount > 0, "amount can not be zero");
        _totalSupply += amount;
        _balances[msg.sender] += amount;
        emit Transfer(address(0), msg.sender, amount);
    }


    //查询账户余额
    function balanceOf(address account) public view virtual returns (uint256) {
        return _balances[account];
    }

    //转账
    function transfer(address to, uint256 amount) public returns (bool) {
        _transfer(msg.sender, to, amount);
        return true;
    }

    //通用转账
    function _transfer(address from, address to, uint256 amount) internal virtual {
        require(from != address(0), "transfer from address can not be zero");
        require(to != address(0), "transfer to address can not be zero");
        require(amount > 0, "amount can not be zero");

        uint256 fromBalance = _balances[from];
        require(fromBalance >= amount, "from balance not enough");
        unchecked {
            _balances[from] -= amount;
        }
        _balances[to] += amount;

        emit Transfer(from, to, amount);
    }

    //查询账户授权额度
    function allowance(address owner, address spender) public view virtual returns (uint256) {
        return _allowances[owner][spender];
    }

    //授权额度，调用者授权给被授权spender的amount额度
    function approve(address spender, uint256 amount) public virtual returns (bool) {
        _approve(msg.sender, spender, amount);
        return true;
    }

    //通用授权额度，调用者授权给被授权spender的amount额度
    function _approve(address owner, address spender, uint256 amount) internal virtual {
        require(owner != address(0), "owner address can not be zero");
        require(spender != address(0), "spender address can not be zero");
        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    //代扣转账，调用者是被授权人，调用此方法，从from账户给to账户转账
    function transferFrom(address from, address to, uint256 amount) public virtual returns (bool) {
        //转账
        _transfer(from, to, amount);

        //授权账户金额变更
        uint256 currentAllowance = _allowances[from][msg.sender];
        require(currentAllowance >= amount, "allowance balance not enough");
        unchecked {
            _approve(from, msg.sender, currentAllowance - amount);
        }

        return true;
    }

}