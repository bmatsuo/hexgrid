package hexgrid
/* 
*  File: coords_test.go
*  Author: Bryan Matsuo [bmatsuo@soe.ucsc.edu] 
*  Created: Sat Jul  2 02:01:22 PDT 2011
*/
import (
    "testing"
)

func TestHighColumn(T *testing.T) {
    if columnIsHigh(0) {
        T.Error("Zero column is high. Zero column should be low.")
    }
    if !columnIsHigh(3) {
        T.Error("3rd column is low. 3rd column should be high.")
    }
    if !columnIsHigh(-3) {
        T.Error("-3rd column is low. -3rd column should be high.")
    }
    if columnIsHigh(6) {
        T.Error("6th column is high. 6th column should be low.")
    }
    if columnIsHigh(-6) {
        T.Error("-6th column is high. -6th column should be low.")
    }
}
