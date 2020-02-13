pragma solidity ^0.4.24;

contract CredentialRepository {
  event Saved(bytes32 indexed id, address by, bytes32 hash);
  
  address owner;

  // id => credential hash
  mapping (bytes32 => bytes32) public credentials;
  
  modifier onlyOwner {
        require(
            msg.sender == owner,
            "Only owner can call this function."
        );
        _;
    }
  
  constructor() public{
      owner = msg.sender;
  }

  function register(bytes32 id, bytes32 hash) public onlyOwner {
    credentials[id] = hash;
    emit Saved(id, msg.sender, hash);
  }

}