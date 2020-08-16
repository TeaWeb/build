<?php

// test command storage

// open access log file
$fp = fopen("/tmp/teaweb-command-storage.log", "a+");

// read access logs from stdin
$stdin = fopen("php://stdin", "r");
while(true) {
    if (feof($stdin)) {
        break;
    }
    $line = fgets($stdin);

    // write to access log file
    fwrite($fp, $line);
}

// close file pointers
fclose($fp);
fclose($stdin);

?>