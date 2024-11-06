i should write about evaluating and analysis
start with a introduction about this chapter 
- my code is writen in golang 
- for zk-snark i use gnark library from go
- the system which is running on is a macpro with m1pro chip
- i analyse my code as component base, which means i divided my code to component like poll creation, vote casting, vote counting, poll mining, vote mining 
- the resone for is, analysing system on blockchain is not easy, because it is on distributed system, so i need to analyse each component to understand the system
- i used random data for my analysis, because i dont have real data. in blow i will explain how i generate random data
-- for poll creation and mining i use the folowing data
--- possible participant size is between []int{1, 2, 5, 10, 15, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000, 100000, 1000000}
--- the question for poll is generated as random between 50 o 150 charachter
--- the options for poll is generated as random between 2 to 5 options
--- the duration is between 1 and 432 

-- poll creation is 5 times Executed with number of polls between 1,5,10,25,50,100,250,500,1000,50000

