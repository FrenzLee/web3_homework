// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Voting {

    mapping (string => uint) public votesRecieved; // Candidate => vote count

    function vote(string memory _candidate) public {
        votesRecieved[_candidate] += 1;
    }

    function getVotes(string memory _candidate) public view returns (uint) {
        return votesRecieved[_candidate];
    }

    function resetVotes() public {
         string[] memory candidatesList = getCandidates();

         for (uint i = 0; i < candidatesList.length; i++) {
            string memory key = candidatesList[i];
            votesRecieved[key] = 0;
         }
    }

    //获取所有候选人名字
    function getCandidates() private pure returns (string[] memory){
        string[] memory candidatesList = new string[](3);
        candidatesList[0] = "candidateA";
        candidatesList[1] = "candidateB";
        candidatesList[2] = "candidateC";

        return candidatesList;
    }



}