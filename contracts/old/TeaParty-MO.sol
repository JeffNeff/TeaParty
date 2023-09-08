// deployed at 0xb74e04F094b7CB084657d61d807b984aD1BbE12A
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// import from node_modules @openzeppelin/contracts v4.0
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/** 
  *@title TeaParty contract
*/
contract TeaParty is ERC20, Ownable, ReentrancyGuard {
    uint256 public _totalSupply;
    // create a mpping of address to uint
    mapping(address => uint) public _teaPartyTransactions;

    
    // Sample constructor
    constructor() ERC20("TeaParty", "TP") {
      _mint(msg.sender, 1000000*(10**uint256(decimals())));
    }
    
    /**
      * @param account (type address) address of recipient
      * @param amount (type uint256) amount of token
      * @dev function use to mint token
    */
    function mint(address account, uint256 amount) public onlyOwner returns (bool sucess) {
      require(account != address(0) && amount != uint256(0), "ERC20: function mint invalid input");
      _totalSupply += amount;
      _mint(account, amount);
      return true;
    }


     /**
     * @dev See {IERC20-totalSupply}.
     */
    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    /** 
      * @dev function to buy token with ether
    */
    function buy() public payable nonReentrant returns (bool sucess) {
      require(msg.sender.balance >= msg.value && msg.value != 0 ether, "buy: function buy invalid input");
      uint256 amount = msg.value + 5;
      _transfer(owner(), _msgSender(), amount);
      _totalSupply += amount;
      return true;
    }


    /** 
      * @param amount (type uint256) amount of ether
      * @dev function use to withdraw ether from contract
    */
    function withdraw(uint256 amount) public onlyOwner returns (bool success) {
      require(amount <= address(this).balance, "withdraw: function withdraw invalid input");
      payable(_msgSender()).transfer(amount);
      return true;
    }

    function purchaseTransaction() public returns (bool success){
      require(amount <= address(this).balance, "purchase: function withdraw invalid input");
      // create and return a random 
      uint test = rnd();
      // add the random number to the mapping
      _teaPartyTransactions[msg.sender] = test;
    }

    function getTeaTransactions() public view returns (uint) {
      return _teaPartyTransactions[msg.sender];
    }

    // remove a specific transaction from the mapping
    function removeTeaTransaction(address transaction) public onlyOwner returns (bool success) {
      delete _teaPartyTransactions[transaction];
      return true;
    }

    function rnd() returns (uint test) {
    uint _test = uint(keccak256("wow"));
    return _test;
  }

}
}