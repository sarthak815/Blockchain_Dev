//SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;
import "@openzeppelin/contracts/utils/Strings.sol";
contract SmartContract{

    uint public hours_time;
    uint public minutes_time;
    string public str = "The number of seconds are: ";
    string public time;
    uint public seconds_time = 0;

    constructor(uint hours_con, uint minutes_con){
        hours_time = hours_con;
        minutes_time = minutes_con;

    }

    function setValues (uint hours_mem, uint minutes_mem) public{
        hours_time = hours_mem;
        minutes_time = minutes_mem;
        seconds_time = hours_time*60*60+minutes_time*60;
        time = Strings.toString(seconds_time);
    }

    function getSeconds() public view returns(string memory){
        
        
        return string(abi.encodePacked(str,time));
    }
}
