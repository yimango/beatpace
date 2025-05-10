from typing import List

class Solution:
    def minSum(self, nums1: List[int], nums2: List[int]) -> int:
        zero1 = 0
        zero2 = 0
        total1 = total2 = 0
        for i in nums1:
            if i == 0:
                zero1 += 1
            total1 += i
        
        for i in nums2:
            if i == 0:
                zero2 += 1
            total2 += i

        if zero2 == 0 and total1 + zero1 > total1:
            return -1
        elif zero1 == 0 and total1 + zero2 > total2:
            return -1
        
        return max(total1 + zero1, total2 + zero2)
    
print(Solution().minSum([3,2,0,1,0], [6,5,0]))  # Output: 6