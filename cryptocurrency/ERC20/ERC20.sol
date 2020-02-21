pragma solidity ^0.6.3;

abstract contract ERC20 {
    function name() virtual public view returns (string memory);
    function symbol() virtual public view returns (string memory);
    function decimals() virtual public view returns (uint);

    function totalSupply() virtual public view returns (uint256);
    function balanceOf(address _owner) virtual public view returns (uint256 balance);
    function approve(address _spender, uint256 _value) virtual public returns (bool success);
    function allowance(address _owner, address _spender) virtual public view returns (uint256 remaining);

    function transfer(address _to, uint256 _value) virtual public returns (bool success);
    function transferFrom(address _from, address _to, uint256 _value) virtual public returns (bool success);

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
 }