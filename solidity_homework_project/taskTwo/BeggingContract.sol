// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import "@openzeppelin/contracts/access/Ownable.sol";

/*
作业3：编写一个讨饭合约
任务目标 
1. 使用 Solidity 编写一个合约，允许用户向合约地址发送以太币。
2. 记录每个捐赠者的地址和捐赠金额。
3. 允许合约所有者提取所有捐赠的资金。

任务步骤 
1. 编写合约
  ○ 创建一个名为 BeggingContract 的合约。
  ○ 合约应包含以下功能：
  ○ 一个 mapping 来记录每个捐赠者的捐赠金额。
  ○ 一个 donate 函数，允许用户向合约发送以太币，并记录捐赠信息。
  ○ 一个 withdraw 函数，允许合约所有者提取所有资金。
  ○ 一个 getDonation 函数，允许查询某个地址的捐赠金额。
  ○ 使用 payable 修饰符和 address.transfer 实现支付和提款。
2. 部署合约
  ○ 在 Remix IDE 中编译合约。
  ○ 部署合约到 Goerli 或 Sepolia 测试网。
3. 测试合约
  ○ 使用 MetaMask 向合约发送以太币，测试 donate 功能。
  ○ 调用 withdraw 函数，测试合约所有者是否可以提取资金。
  ○ 调用 getDonation 函数，查询某个地址的捐赠金额。

任务要求 
1. 合约代码：
  ○ 使用 mapping 记录捐赠者的地址和金额。
  ○ 使用 payable 修饰符实现 donate 和 withdraw 函数。
  ○ 使用 onlyOwner 修饰符限制 withdraw 函数只能由合约所有者调用。
2. 测试网部署：
  ○ 合约必须部署到 Goerli 或 Sepolia 测试网。
3. 功能测试：
  ○ 确保 donate、withdraw 和 getDonation 函数正常工作。

提交内容 
1. 合约代码：提交 Solidity 合约文件（如 BeggingContract.sol）。
2. 合约地址：提交部署到测试网的合约地址。
3. 测试截图：提交在 Remix 或 Etherscan 上测试合约的截图。

额外挑战（可选） 
1. 捐赠事件：添加 Donation 事件，记录每次捐赠的地址和金额。（已开发）
2. 捐赠排行榜：实现一个功能，显示捐赠金额最多的前 3 个地址。（已开发）
3. 时间限制：添加一个时间限制，只有在特定时间段内才能捐赠

合约创建测试区块链地址：https://sepolia.etherscan.io/tx/0x3d56abb1c26bb16898617908348f9aad009fd6d5ab278a832fad228e3edca64c
*/

contract BeggingContract is Ownable {
    //记录每个捐赠者的捐赠金额
    mapping(address => uint256) private _balances;
    // 用于防止重入攻击的互斥锁
    bool private locked;
    //记录排名前三位捐赠者的地址
    address[3] private topDonors;

    //定义事件
    event Donation(address indexed from, uint256 value);

    //互斥锁，防止重入攻击
    modifier noReentrancy() {
        require(!locked, "No reentrancy allowed!");
        locked = true;
        _;
        locked = false;
    }

    constructor() Ownable(msg.sender) {}

    //用户捐赠
    function donate() external payable {
        require(msg.value > 0, "donate value must > 0");
        _balances[msg.sender] += msg.value;
        updateDonors(msg.sender);
        emit Donation(msg.sender, msg.value);
    }

    //合约所有者提取所有资金
    function withdraw() external onlyOwner noReentrancy {
        //这个合约中的所有代币
        uint256 balance = address(this).balance;
        require(balance > 0, "No balance to withdraw");
        //转账
        payable(msg.sender).transfer(balance);
    }

    //查询某个地址的捐赠金额
    function getDonation(address addr) external view returns(uint256) {
        return _balances[addr];
    }

    //获取前三名捐赠者
    function getTopDonors() external view returns (address[3] memory, uint256[3] memory) {
        address[3] memory donors;
        uint256[3] memory amounts;
        for (uint i = 0; i < topDonors.length; i++) {
            donors[i] = topDonors[i];
            amounts[i] = _balances[topDonors[i]];
        }
        return (donors, amounts);
    }

    //更新topDonors
    function updateDonors(address donorAddr) internal {

        if (topDonors.length == 0) {
            topDonors[0] = donorAddr;
        } else if (topDonors.length == 1) {
            topDonors[1] = donorAddr;
            sortedDonors();
        } else if (topDonors.length == 2) {
            topDonors[2] = donorAddr;
            sortedDonors();
        } else if (topDonors.length == 3) {
            uint256 donorBalance = _balances[donorAddr];
            if (donorBalance >= _balances[topDonors[2]]) {
                topDonors[2] = donorAddr;
                sortedDonors();
            }
        }
    }

    //排序
    function sortedDonors() internal {
        for (uint i = 0; i < topDonors.length; i++){
            for (uint j = i + 1; j < topDonors.length; j++) {
                if (_balances[topDonors[i]] < _balances[topDonors[j]]) {
                    address temp = topDonors[i];
                    topDonors[i] = topDonors[j];
                    topDonors[j] = temp;
                }
            }
        }
       
        
    }
}