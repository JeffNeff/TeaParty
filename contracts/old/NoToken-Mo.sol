// deployed at 0x96168158d0F8085cb2A84befeE4a8a94Ed34f79a
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// import from node_modules @openzeppelin/contracts v4.0
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/** 
  *@title TeaParty contract
*/
contract TeaParty is Ownable, ReentrancyGuard {
    // create a mpping of address to uint
    mapping(address => uint) public _teaPartyTransactions;

    /** 
      * @param amount (type uint256) amount of ether
      * @dev function use to withdraw ether from contract
    */
    function withdraw(uint256 amount) public onlyOwner returns (bool success) {
      require(amount <= address(this).balance, "withdraw: function withdraw invalid input");
      payable(_msgSender()).transfer(amount);
      return true;
    }

    function purchaseTransaction() public payable returns (bool success){
      require(msg.sender.balance >= msg.value && msg.value != 0 ether, "purchaseTransaction: function buy invalid input");
      // create and return a random 
      uint test = rnd();
      // add the random number to the mapping
      _teaPartyTransactions[msg.sender] = test;
      return true;
    }


    function getTeaTransactions() public view returns (uint) {
      return _teaPartyTransactions[msg.sender];
    }

    // remove a specific transaction from the mapping
    function removeTeaTransaction(address transaction) public onlyOwner returns (bool success) {
      delete _teaPartyTransactions[transaction];
      return true;
    }

    function rnd() private returns (uint test) {
    uint _test = uint(keccak256("wow"));
    return _test;
  }

}