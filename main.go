//go:generate swagger generate spec

// TODO: currently no support for piped mode.

package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	ce "github.com/engelch/go_libs/v2"
	cli "github.com/urfave/cli/v2"
)

const appVersion = "0.0.5"
const appName = "sha256base64"

// These CLI options are used more than once below. So let's use constants that we do not get
// misbehaviour by typoos.
const _debug = "debug" // long (normal) name of CLI option
const _v2 = "v2"       // dito

// =======================================================================================
// checkOptions checks the command line options if properly set or in range.
// POST: exactly one keyfile is not mt.
func checkOptions(c *cli.Context) error {
	if c.Bool(_debug) {
		ce.CondDebugSet(true)
	}
	ce.CondDebugln("Debug is enabled.")
	return nil
}

// commandLineOptions just separates the definition of command line options ==> creating a shorter main
func commandLineOptions() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    _debug,
			Aliases: []string{"d"},
			Value:   false,
			Usage:   "OPTIONAL: enable debug",
		},
	}
}

// It creates for a textfile
//
// -rw-r--r--  1 engelch  staff    32 Mar 15 16:42 text.sha256bin
// -rw-r--r--  1 engelch  staff    64 Mar 15 16:42 text.sha256
// -rw-r--r--  1 engelch  staff    44 Mar 15 16:42 text.sha256b64
// -rw-r--r--  1 engelch  staff   256 Mar 15 16:42 text.sha256bytes
//
// text.sha256 ---------------
// 322dfb6889f945ed825565882245209261fc354e1bf5d64c193bb96240f9c33e
// text.sha256b64 ---------------
// Mi37aIn5Re2CVWWIIkUgkmH8NU4b9dZMGTu5YkD5wz4=
// text.sha256bin ---------------
// 2-?h??E?Ue?"E ?a?5N?L;?b@??>
// text.sha256bytes ---------------
// Mi37aIn5Re2CVWWIIkUgkmH8NU4b9dZMGTu5YkD5wz4=

// main start routine
func main() {
	app := cli.NewApp() // global var, see discussion above
	app.Flags = commandLineOptions()
	app.Name = appName
	app.Version = appVersion
	app.Usage = "It calculates the SHA-256 sum for the 1st argument supplied. This information is stored in different forms.\n" +
		"             This information is required as openssl with PSS padding (v2) only supports binary hashes of 32 bytes size.\n" +
		"             Without PSS padding (v1), openssl is more tolerant."

	app.Action = func(c *cli.Context) error {
		err := checkOptions(c)
		ce.ExitIfError(err, 9, "checkOptions")

		bytes, err := ioutil.ReadFile(c.Args().Get(0))
		ce.ExitIfError(err, 1, "reading file")
		extension := filepath.Ext(c.Args().Get(0))
		basename := c.Args().Get(0)[0 : len(c.Args().Get(0))-len(extension)]

		sha256 := ce.Sha256bytes2bytes(bytes)

		// binary output
		err = ioutil.WriteFile(basename+".sha256bin", sha256, 0644)
		ce.ExitIfError(err, 2, "write sha256bin")
		fmt.Printf("SHA256 binary:             %v\n", sha256)

		// hex output
		err = ioutil.WriteFile(basename+".sha256", []byte(fmt.Sprintf("%x", sha256)), 0644)
		ce.ExitIfError(err, 3, "write sha256")
		fmt.Printf("SHA256 hex (normal):       %x\n", sha256)

		// b64 output
		err = ioutil.WriteFile(basename+".sha256b64", []byte(base64.StdEncoding.EncodeToString(sha256)), 0644)
		ce.ExitIfError(err, 4, "write sha256b64")
		fmt.Printf("SHA256 b64-coded (string): %v\n", []byte(base64.StdEncoding.EncodeToString(sha256)))
		fmt.Printf("SHA256 b64-coded (string): %s\n", base64.StdEncoding.EncodeToString(sha256))

		// same output as b64
		var outbytes [256]byte
		base64.StdEncoding.Encode(outbytes[:], sha256)
		err = ioutil.WriteFile(basename+".sha256bytes", outbytes[:], 0644)
		ce.ExitIfError(err, 5, "write sha256bytes")
		fmt.Printf("SHA256 b64-coded (bytes)   %v\n", outbytes[:])

		return nil
	}
	_ = app.Run(os.Args)
}

// eof
