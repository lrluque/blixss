## Blind XSS Tool (BLiXSS)

BLiXSS is a command-line tool designed for blind XSS (Cross-Site Scripting) testing. It crafts a malicious payload to inject into web application parameters, allowing you to detect potential vulnerabilities.

![blixss](https://github.com/lrluque/blixss/assets/16742563/c7ebe0a2-2aef-4662-b29a-f0260f4cb2e8)



### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/blixss.git
```

Navigate to the `blixss` directory:

```bash
cd blixss
```

Build the executable:

```bash
go build blixss.go
```

### Usage

BLiXSS requires the following parameters:

- `-t`: Target URL (e.g., `http://example.com`)
- `-b`: Body strings with the parameters of the request (e.g., `"parameter1=XSS&parameter2=test2&parameter3=XSS"`)
- `-l`: URL to forward the requests to (e.g., `http://10.10.15.122:45000`)
- `-d`: Specifies a custom directory to make the GET request. If not specified, it will attach `/<<paramName>>` on the request.

Example usage:

```bash
./blixss -t "http://example.com" -b "parameter1=XSS&parameter2=test2&parameter3=XSS" -l "http://10.10.15.122:45000" -d "custom/request/directory"
```

Parameter values different from 'XSS' will not be tested.

### Disclaimer

This tool is for educational purposes only. Do not use it for any illegal activities. I am not responsible for any misuse or damage caused by BLiXSS.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
