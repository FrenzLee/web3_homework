// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import "forge-std/console.sol";
import {OpMathCalV2} from "../src/OpMathCalV2.sol";

contract OpMathCalV1Test is Test {
    OpMathCalV2 public opMathCalV2;

    // 记录V2版本的 Gas 消耗
    uint256 public addGasV2;
    uint256 public subGasV2;
    uint256 public mulGasV2;
    uint256 public divGasV2;
    uint256 public resetGasV2;
    uint256 public squareGasV2;
    uint256 public cubeGasV2;

    function setUp() public {
        opMathCalV2 = new OpMathCalV2();
    }

    function test_AddV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.add(10, 20);
        uint256 gasEnd = gasleft();
        addGasV2 = gasStart - gasEnd;

        assertEq(result, 30, "Addition result incorrect");
        assertEq(opMathCalV2.getLastOperation(), 0, "Last operation code incorrect");
        assertEq(opMathCalV2.getResult(), 30, "Stored result incorrect");

        console.log("V2 Add operation gas used:", addGasV2);
    }

    function test_SubtractV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.sub(50, 20);
        uint256 gasEnd = gasleft();
        subGasV2 = gasStart - gasEnd;

        assertEq(result, 30, "Subtraction result incorrect");
        console.log("V2 Subtract operation gas used:", subGasV2);
    }

    function test_MultiplyV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.mul(5, 6);
        uint256 gasEnd = gasleft();
        mulGasV2 = gasStart - gasEnd;

        assertEq(result, 30, "Multiplication result incorrect");
        console.log("V2 Multiply operation gas used:", mulGasV2);
    }

    function test_DivideV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.div(60, 2);
        uint256 gasEnd = gasleft();
        divGasV2 = gasStart - gasEnd;

        assertEq(result, 30, "Division result incorrect");
        console.log("V2 Divide operation gas used:", divGasV2);
    }

    function test_SquareV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.square(5);
        uint256 gasEnd = gasleft();
        squareGasV2 = gasStart - gasEnd;

        assertEq(result, 25, "Square result incorrect");
        console.log("V2 Square operation gas used:", squareGasV2);
    }

    function test_CubeV2() public {
        uint256 gasStart = gasleft();
        uint256 result = opMathCalV2.cube(3);
        uint256 gasEnd = gasleft();
        cubeGasV2 = gasStart - gasEnd;

        assertEq(result, 27, "Cube result incorrect");
        console.log("V2 Cube operation gas used:", cubeGasV2);
    }

    function test_ResetV2() public {
        opMathCalV2.add(10, 20); // Set a value first

        uint256 gasStart = gasleft();
        opMathCalV2.reset();
        uint256 gasEnd = gasleft();
        resetGasV2 = gasStart - gasEnd;

        assertEq(opMathCalV2.getResult(), 0, "Reset result incorrect");
        console.log("V2 Reset operation gas used:", resetGasV2);
    }

    function test_GasSummaryV2() public {
        // 执行所有操作并记录 Gas 消耗
        opMathCalV2.add(10, 20);
        opMathCalV2.sub(30, 10);
        opMathCalV2.mul(5, 4);
        opMathCalV2.div(40, 2);
        opMathCalV2.square(6);
        opMathCalV2.cube(3);
        opMathCalV2.reset();

        console.log("=== V2 Gas Consumption Summary ===");
        console.log("V2 Add gas used:", addGasV2);
        console.log("V2 Subtract gas used:", subGasV2);
        console.log("V2 Multiply gas used:", mulGasV2);
        console.log("V2 Divide gas used:", divGasV2);
        console.log("V2 Square gas used:", squareGasV2);
        console.log("V2 Cube gas used:", cubeGasV2);
        console.log("V2 Reset gas used:", resetGasV2);
    }

}