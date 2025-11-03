// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";
import {OpMathCalV1} from "../src/OpMathCalV1.sol";

contract OpMathCalV1Test is Test {
    OpMathCalV1 public opMathCalV1;

    // 记录优化版本的 Gas 消耗
    uint256 public addGasV1;
    uint256 public subGasV1;
    uint256 public mulGasV1;
    uint256 public divGasV1;
    uint256 public resetGasV1;

    function setUp() public {
        opMathCalV1 = new OpMathCalV1();
    }

    function test_AddV1() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV1.add(10, 20);
        uint256 gasEnd = gasleft();
        addGasV1 = gasStart - gasEnd;

        assertEq(result, 30, "Addition result incorrect");
        (uint128 a, uint128 b) = opMathCalV1.getLastCalculation();
        assertEq(a, 10, "Last operand A incorrect");
        assertEq(b, 20, "Last operand B incorrect");

        console.log("V1 Add operation gas used:", addGasV1);
    }

    function test_SubtractV1() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV1.sub(50, 20);
        uint256 gasEnd = gasleft();
        subGasV1 = gasStart - gasEnd;

        assertEq(result, 30, "Subtraction result incorrect");
        console.log("V1 Subtract operation gas used:", subGasV1);
    }

    function test_MultiplyV1() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV1.mul(5, 6);
        uint256 gasEnd = gasleft();
        mulGasV1 = gasStart - gasEnd;

        assertEq(result, 30, "Multiplication result incorrect");
        console.log("V1 Multiply operation gas used:", mulGasV1);
    }

    function test_DivideV1() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV1.div(60, 2);
        uint256 gasEnd = gasleft();
        divGasV1 = gasStart - gasEnd;

        assertEq(result, 30, "Division result incorrect");
        console.log("V1 Divide operation gas used:", divGasV1);
    }

    function test_ResetV1() public {
        opMathCalV1.add(10, 20); // Set a value first

        uint256 gasStart = gasleft();
        opMathCalV1.reset();
        uint256 gasEnd = gasleft();
        resetGasV1 = gasStart - gasEnd;

        assertEq(opMathCalV1.getResult(), 0, "Reset result incorrect");
        console.log("V1 Reset operation gas used:", resetGasV1);
    }

    function test_GasSummaryV1() public {
        // 执行所有操作并记录 Gas 消耗
        opMathCalV1.add(10, 20);
        opMathCalV1.sub(30, 10);
        opMathCalV1.mul(5, 4);
        opMathCalV1.div(40, 2);
        opMathCalV1.reset();

        console.log("=== V1 Gas Consumption Summary ===");
        console.log("V1 Add gas used:", addGasV1);
        console.log("V1 Subtract gas used:", subGasV1);
        console.log("V1 Multiply gas used:", mulGasV1);
        console.log("V1 Divide gas used:", divGasV1);
        console.log("V1 Reset gas used:", resetGasV1);
    }

}