// deployed at 0xDAD416F84E67d4B37c12c37979DC4E8d07Fc83d2
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

contract TeaParty is Ownable, ReentrancyGuard {
    using EnumerableSet for EnumerableSet.AddressSet;
    EnumerableSet.AddressSet private _holders;

    function addHolder() public {
        EnumerableSet.add(_holders, msg.sender);
    }

    function getHolders() public view returns (address[] memory) {
        uint256 length = EnumerableSet.length(_holders);
        address[] memory holders = new address[](length);
        for (uint256 i = 0; i < length; i++) {
            holders[i] = EnumerableSet.at(_holders, i);
        }
        return holders;
    }

    mapping(address => uint) private _teaPartyTransactions;
    uint private _openTransactions;

    event TransactionCreated(address indexed from, uint transactionId);

    function withdraw(uint256 amount) public onlyOwner {
        require(amount <= address(this).balance, "TeaParty: withdraw amount exceeds balance");
        payable(owner()).transfer(amount);
    }

    function deposit() public payable {
        require(msg.value > 0, "TeaParty: deposit amount must be greater than zero");
    }

    function createTransaction() public payable nonReentrant returns (uint) {
        require(msg.value > 0, "TeaParty: transaction amount must be greater than zero");
        _teaPartyTransactions[msg.sender] = _openTransactions;
        EnumerableSet.add(_holders, msg.sender);
        _openTransactions++;
        emit TransactionCreated(msg.sender, _openTransactions);
        return _openTransactions;
    }

    function getTransaction(address participant) public view returns (uint) {
        require(_teaPartyTransactions[participant] != 0, "TeaParty: transaction does not exist");
        return _teaPartyTransactions[participant];
    }

    function removeTransaction(address participant) public onlyOwner {
        delete _teaPartyTransactions[participant];
        _openTransactions--;
    }

    function rnd() private view returns (uint) {
        return uint(keccak256(abi.encode(block.timestamp, block.prevrandao, msg.sender)));
    }

    function openTransactions() public view returns (uint) {
        return _openTransactions;
    }

    function getTeaPartyTransactions() public view onlyOwner returns (address[] memory, uint[] memory) {
        address[] memory keys = new address[](_openTransactions);
        uint[] memory values = new uint[](_openTransactions);

        for (uint i = 0; i < _openTransactions; i++) {
            keys[i] = getAddressAtIndex(i);
            values[i] = _teaPartyTransactions[keys[i]];
        }

        return (keys, values);
    }

    function getAddressAtIndex(uint index) private view returns (address) {
        require(index < _openTransactions, "TeaParty: index out of range");
        return EnumerableSet.at(_holders, index);
    }
}