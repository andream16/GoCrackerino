# GoCrackerino
##Golang Key Cracker Which Implements Meet In the Middle Strategy.

Given a known `cipher-text` and a known `plain-text`, it finds two keys `k1` and `k2` in a way that `Ek1(plain) == Dk2(cipher)`
where `E` is the encryption function and `D` is the decryption function.

It tries all the possible keys in the range of `24/28/32 bits` and prints which are the keys to be used for the encryption and decryption.

Keys can only be `6/7/8 Digits HEX` and `Cipher/Plain` can only be `16 Digits HEX`.

This cracker was written in `Golang` since I needed to run both `Encryption` and `Decryption` in parallel and `goroutines` made it easy.

#What you need to run the project:

    - Go   : ~1.7
    - Ocaml: ~4.03.0
    - AWK  : ~4.0.1
    - Grep : ~2.x
    
Download the project, extract the contents of `sp-network for challenge mitm.zip`, then run `make` on its folder, two `executables` and more files will be generated. Move `encrypt` and `decrypt` executables to the main project's folder or just make sure to specify their paths on `go-crackerino.go`.

You can run `encryption` using: `./encrypt -s key cipherText` and `decryption` using `./decrypt -s key plainText`.

    - 24-bit:   000000 -   FFFFFF
    - 28-bit:  0000000 -  FFFFFFF
    - 32-bit: 00000000 - FFFFFFFF

`Plain/Cipher` couples to be used to test the application:

    - 24-bit: C330C9CBD01DFBA0 - E10C65124518DB05
    - 28-bit: 492A5BB83F8A3F95 - 47090A6AC4A56798
    - 32-bit: FFC7C9E5694ABFF7 - 98AC59F25448FFAC
    
Make sure to have `~20GB` of free space for the 32-bit challenge since it will write huge amounts of couples `key cipher` in a file.

When `Ek1(C) == Dk2(P)` you found the right keys `k1, k2` to read all encrypted messages.
