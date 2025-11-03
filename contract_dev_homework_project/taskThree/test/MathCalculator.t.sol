// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";
import {MathCalculator} from "../src/MathCalculator.sol";

contract MathCalculatorTest is Test {
    MathCalculator public calculator;

    uint256 public addGas;
    uint256 public subGas;
    uint256 public mulGas;
    uint256 public divGas;
    uint256 public resetGas;

    function setUp() public {
        calculator = new MathCalculator();
    }

    function test_Add() public {
        uint256 gasStart = gasleft();
        uint256 result = calculator.add(10, 20);
        uint256 gasEnd = gasleft();
        addGas = gasStart - gasEnd;

        assertEq(result, 30, "Addition result incorrect");
        assertEq(calculator.getResult(), 30, "Stored result incorrect");

        console.log("Add operation gas used:", addGas);
    }

    function test_Subtract() public {
        calculator.add(50, 20); // Set initial value

        uint256 gasStart = gasleft();
        uint256 result = calculator.sub(50, 20);
        uint256 gasEnd = gasleft();
        subGas = gasStart - gasEnd;

        assertEq(result, 30, "Subtraction result incorrect");
        assertEq(calculator.getResult(), 30, "Stored result incorrect");

        console.log("Subtract operation gas used:", subGas);
    }

    function test_Multiply() public {
        uint256 gasStart = gasleft();
        uint256 result = calculator.mul(5, 6);
        uint256 gasEnd = gasleft();
        mulGas = gasStart - gasEnd;

        assertEq(result, 30, "Multiplication result incorrect");
        assertEq(calculator.getResult(), 30, "Stored result incorrect");

        console.log("Multiply operation gas used:", mulGas);
    }

    function test_Divide() public {
        uint256 gasStart = gasleft();
        uint256 result = calculator.div(60, 2);
        uint256 gasEnd = gasleft();
        divGas = gasStart - gasEnd;

        assertEq(result, 30, "Division result incorrect");
        assertEq(calculator.getResult(), 30, "Stored result incorrect");

        console.log("Divide operation gas used:", divGas);
    }

    function test_Reset() public {
        calculator.add(10, 20); // Set a value first

        uint256 gasStart = gasleft();
        calculator.reset();
        uint256 gasEnd = gasleft();
        resetGas = gasStart - gasEnd;

        assertEq(calculator.getResult(), 0, "Reset result incorrect");

        console.log("Reset operation gas used:", resetGas);
    }

    function testFuzz_AddAssociative(uint256 a, uint256 b, uint256 c) public {
        // 防止溢出
        vm.assume(b <= type(uint256).max - c); 
        vm.assume(a <= type(uint256).max - b - c);

        calculator.add(a, b);
        uint256 result1 = calculator.add(calculator.getResult(), c);

        calculator.reset();
        calculator.add(b, c);
        uint256 result2 = calculator.add(a, calculator.getResult());

        assertEq(result1, result2, "Addition associative property failed");
    }



    // Gas 性能汇总测试
    function test_GasSummary() public {
        // 执行所有操作并记录 Gas 消耗
        calculator.add(10, 20);
        calculator.sub(30, 10);
        calculator.mul(5, 4);
        calculator.div(40, 2);
        calculator.reset();

        console.log("=== Gas Consumption Summary ===");
        console.log("Add gas used:", addGas);
        console.log("Subtract gas used:", subGas);
        console.log("Multiply gas used:", mulGas);
        console.log("Divide gas used:", divGas);
        console.log("Reset gas used:", resetGas);
        console.log("Total gas for all operations:",
                    addGas + subGas + mulGas + divGas + resetGas);
    }

}
