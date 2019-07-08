secspy(1) -- Utility to gather data, image, events and videos from SecuritySpy
===

## SYNOPSIS

`secspy -u admin -p password -c cameras`

## DESCRIPTION

* This application provides a command-line interface into the SecuritySpy.

## OPTIONS

`secspy [-u <username>] [-p <password>] [-U <url>] [-c <cmd>] [-a <arg]`

    -u, --user <username>
      Username to authenticate with. Can also be set with env variable SECSPY_USERNAME

    -p, --pass <password>
      Password to authenticate with. Can also be set with env variable SECSPY_PASSWORD

    -U, --url <url>
      SecuritySpy URL. Default is http://127.0.0.1:8000/

    -s, --verify-ssl
      If using SSL (it breaks RTSP), pass this flag to validate the SSL certificate.

    -c, --command
        Command to run. Choices:
        e|events   - Watch (some of) the event stream
        l|callback - Another way to watch event stream
        t|trigger  - Trigger motion on a camera
        p|picture  - Save a live picture to a file
        v|video    - Save a live 10 second video to a file
        c|cameras  - List all cameras and data
        f|files    - Show saved media files for a camera
        d|download - Downloads a saved media file
        z|ptz      - Controls a PTZ camera

    -a, --arg <arg>
        Some commands require an argument. Use this to provide the arg.

    -h, --help
        Display usage and exit.

## AUTHOR

* David Newhall II - January 2019

## LOCATION

* https://github.com/davidnewhall/SecSPyCLI
* /usr/local/bin/secspy
