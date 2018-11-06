pragma solidity ^0.4.0;

contract SimpleStorage {
    uint storedData;
    event getX(uint,uint);

    function SimpleStorage() public{
        storedData = 5;
    }

    function set(uint x) public {
        storedData = x;
    }

    function get() public returns (uint) {
        emit getX(1,2);
        return storedData;
    }
}