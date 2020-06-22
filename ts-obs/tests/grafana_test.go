package tests

import (
	"net"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

func testGrafanaPortForward(t testing.TB, port string) {
	var portforward *exec.Cmd

	if port == "" {
		t.Logf("Running 'ts-obs grafana port-forward'")
		portforward = exec.Command("ts-obs", "grafana", "port-forward")
	} else {
		t.Logf("Running 'ts-obs grafana port-forward -p %v'\n", port)
		portforward = exec.Command("ts-obs", "grafana", "port-forward", "-p", port)
	}

	err := portforward.Start()
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(4 * time.Second)

	if port == "" {
		port = "8080"
	}

	_, err = net.DialTimeout("tcp", "localhost:"+port, time.Second)
	if err != nil {
		t.Fatal(err)
	}

	portforward.Process.Signal(syscall.SIGINT)
}

func testGrafanaGetPass(t testing.TB) {
	var getpass *exec.Cmd

	t.Logf("Running 'ts-obs grafana get-initial-password'")
	getpass = exec.Command("ts-obs", "grafana", "get-initial-password")

	out, err := getpass.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}

func testGrafanaChangePass(t testing.TB, newpass string) {
	var changepass *exec.Cmd

	t.Logf("Running 'ts-obs grafana change-password %v'\n", newpass)
	changepass = exec.Command("ts-obs", "grafana", "change-password", newpass)

	out, err := changepass.CombinedOutput()
	if err != nil {
		t.Logf(string(out))
		t.Fatal(err)
	}
}

func TestGrafana(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Grafana tests")
	}

	testGrafanaPortForward(t, "")
	testGrafanaPortForward(t, "1235")
	testGrafanaPortForward(t, "2348")
	testGrafanaPortForward(t, "7390")

	testGrafanaGetPass(t)
	testGrafanaChangePass(t, "kraken")
	testGrafanaChangePass(t, "cereal")
	testGrafanaChangePass(t, "23498MSDF(*9389m*(@#M24309mDj")
	testGrafanaGetPass(t)
}
