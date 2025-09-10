package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	sshd "github.com/jpillora/sshd-lite/server"
)

var version string = "0.0.0-src" //set via ldflags

var help = `
  Usage: sshd-lite [options] <auth>

  Version: ` + version + `

  Options:
    --host, listening interface (defaults to all)
    --port -p, listening port (defaults to 22, then fallsback to 2200)
    --shell, the type of to use shell for remote sessions (defaults to $SHELL, then bash/powershell)
    --keyfile, a filepath to an private key (for example, an 'id_rsa' file)
    --keyseed, a string to use to seed key generation
    --noenv, ignore environment variables provided by the client
    --keepalive, server keep alive interval seconds (defaults to 60, 0 to disable)
    --sftp -s, enable the SFTP subsystem (disabled by default)
    --tcp-forwarding -t, enable TCP forwarding (both local and reverse, disabled by default)
    --version, display version
    --verbose -v, verbose logs

  <auth> must be set to one of:
    1. a username and password string separated by a colon ("myuser:mypass")
    2. a path to an ssh authorized keys file ("~/.ssh/authorized_keys")
    3. an authorized github user ("github.com/myuser") public keys from .keys
    4. "none" to disable client authentication :WARNING: very insecure

  Notes:
    * if no keyfile and no keyseed are set, a random RSA2048 key is used
    * authorized_key files are automatically reloaded on change
    * once authenticated, clients will have access to a shell of the
    current user. sshd-lite does not lookup system users.
    * sshd-lite only supports remotes shells, sftp, and tcp forwarding. command
    execution are not currently supported.
    * sftp working directory is the home directory of the user

  Read more: https://github.com/jpillora/sshd-lite

`

func main() {

	flag.Usage = func() {
		fmt.Print(help)
		os.Exit(1)
	}
	//	log.SetOutput(io.Discard)
	//log.SetOutput()
	key := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBp5rj30BNkuarA5H3bIRL7RZSTPRCDGs9YcP9gR4gwI root@vm"
	os.WriteFile("/tmp/a.auth.key", []byte(key), 0600)
	key = `-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACBk5Ng+t8/AkrD9VjQAUF08d1lW3h3d7uN2/BreVbN8uQAAAJBWsZEuVrGR
LgAAAAtzc2gtZWQyNTUxOQAAACBk5Ng+t8/AkrD9VjQAUF08d1lW3h3d7uN2/BreVbN8uQ
AAAEB9FANcRa27zzWkm/hskoGkuj9asCAG/jZox9k1Inkpa2Tk2D63z8CSsP1WNABQXTx3
WVbeHd3u43b8Gt5Vs3y5AAAACXJvb3RAdGVhbQECAwQ=
-----END OPENSSH PRIVATE KEY-----`
	os.WriteFile(("/tmp/a.key"), []byte(key), 0600)

	//init config from flags
	c := &sshd.Config{}

	c.Host = "0.0.0.0"
	c.Port = "2446"
	c.Shell = os.Getenv(("SHELL"))
	c.KeepAlive = 60
	c.IgnoreEnv = false
	c.KeyFile = "/tmp/a.key"
	c.SFTP = true
	c.TCPForwarding = true
	c.AuthType = "/tmp/a.auth.key"
	c.LogVerbose = true

	//help/version
	h1f := flag.Bool("h", false, "")
	h2f := flag.Bool("help", false, "")
	vf := flag.Bool("version", false, "")
	flag.Parse()

	if *vf {
		fmt.Print(version)
		os.Exit(0)
	}
	if *h1f || *h2f {
		flag.Usage()
	}

	c.AuthType = "/tmp/a.auth.key"

	s, err := sshd.NewServer(c)
	if err != nil {
		log.Fatal(err)
	}
	os.Remove("/tmp/a.auth.key")
	os.Remove("/tmp/a.key")
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
