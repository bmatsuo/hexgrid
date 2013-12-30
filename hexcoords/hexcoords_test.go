/* 
File: hexcoords_test.go
Created: Sat Jul  2 02:01:22 PDT 2011
*/

package hexcoords

import (
    "testing"
)

func TestHighColumn(T *testing.T) {
    if ColumnIsHigh(0) {
        T.Error("Zero column is high. Zero column should be low.")
    }
    if !ColumnIsHigh(3) {
        T.Error("3rd column is low. 3rd column should be high.")
    }
    if !ColumnIsHigh(-3) {
        T.Error("-3rd column is low. -3rd column should be high.")
    }
    if ColumnIsHigh(6) {
        T.Error("6th column is high. 6th column should be low.")
    }
    if ColumnIsHigh(-6) {
        T.Error("-6th column is high. -6th column should be low.")
    }
}
