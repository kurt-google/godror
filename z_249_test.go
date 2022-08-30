package godror_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/godror/godror"
)

func TestOpenCloseCrash(t *testing.T) {
	c1, e1 := sql.Open("godror", "oracle://?sysdba=1")
	t.Logf("e1: %v", e1)
	c2, e2 := sql.Open("godror", "oracle://?sysdba=1")
	t.Logf("e2: %v", e2)

	join := make(chan struct{})
	waitForSch := make(chan struct{})
	go func() {
		waitForSch <- struct{}{}
		e1 := c1.PingContext(context.TODO())
		t.Logf("conn1 done: %v", e1)
		c1.Close()
		join <- struct{}{}
	}()
	go func() {
		<-waitForSch
		e2 := c2.PingContext(context.TODO())
		t.Logf("conn2 done: %v", e2)
		c2.Close()
		join <- struct{}{}
	}()
	<-join
	<-join
}
