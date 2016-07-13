//https://leetcode.com/problems/sum-of-two-integers/
// 371. Sum of Two Integers

  
//Calculate the sum of two integers a and b, but you are not allowed to use the operator + and -.

//Example:
//Given a = 1 and b = 2, return 3.

class Solution {
public:
    int getSum(int x, int y) { 
        
    if (y == 0)
        return x;
    else
        return getSum( x ^ y, (x & y) << 1);
    }
};
