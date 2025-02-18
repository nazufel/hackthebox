# Flag Command 

Solution for the [Flag Command](https://app.hackthebox.com/challenges/Flag%20Command/walkthroughs) web challenge.

## Methodology

1. **Initial Reconnaissance**:
   - Navigated to the target URL provided by the challenge
   - Inspected the page source for any clues or hidden content
   - Found there are multiple JS files in the `<script>` tags under `/static/terminal/js` and started inspecting them
   - Noticed the `main.js` file was making api calls to `api/options` and `api/monitor`
   - Made a request to `/api/options` and got a list of options. 
        ```json
        {
        "allPossibleCommands": {
            "1": [
                "HEAD NORTH",
                "HEAD WEST",
                "HEAD EAST",
                "HEAD SOUTH"
            ],
            "2": [
                "GO DEEPER INTO THE FOREST",
                "FOLLOW A MYSTERIOUS PATH",
                "CLIMB A TREE",
                "TURN BACK"
            ],
            "3": [
                "EXPLORE A CAVE",
                "CROSS A RICKETY BRIDGE",
                "FOLLOW A GLOWING BUTTERFLY",
                "SET UP CAMP"
            ],
            "4": [
                "ENTER A MAGICAL PORTAL",
                "SWIM ACROSS A MYSTERIOUS LAKE",
                "FOLLOW A SINGING SQUIRREL",
                "BUILD A RAFT AND SAIL DOWNSTREAM"
            ],
            "secret": [
                "Blip-blop, in a pickle with a hiccup! Shmiggity-shmack"
            ]
            }
        }
        ```
   - Noticed there was list of lists of commands in the response and there was a `secret` command in the list. 

2. **Exploitation**:
    - Copied the `secret` command and made a request to `/api/monitor` with the following JSON:
        ```json
        {
            "command": "Blip-blop, in a pickle with a hiccup! Shmiggity-shmack"
        }
        ```
    - Got the flag in the response
        ```json
        {
            "message": "HTB{D3v3l0p3r_t00l5_4r3_b35t__t0015_wh4t_d0_y0u_Th1nk??}"
        }
        ```
3. **Proof of Concept**:
    - [main.go](main.go) is a proof of concept that makes the request to `/api/monitor` with the `secret` command and prints the flag.

## Impact
* An unauthenticated attacker can access sensitive API endpoints (`/api/options` and `/api/monitor`)
* The API reveals "secret" commands that appear to be intended to be hidden
* The API allows execution of these commands without any authentication or authorization checks
* The endpoint returns sensitive data (the flag in this case, but in a real system this could be other confidential information)

### Severity: MEDIUM to HIGH

### Reasons:
1. Authentication Bypass: The complete lack of authentication controls is a significant security issue
2. Information Disclosure: The API openly reveals internal commands and sensitive data
3. No Rate Limiting: There appears to be no protection against automated requests/brute force attempts

### Mitigating Factors:
This appears to be a CTF challenge, so the actual data exposed is just the challenge flag
The endpoint seems to only return predefined responses rather than allowing arbitrary command execution

### Real-World Implications:
If this were a real production API:
* Attackers could enumerate valid commands and potentially access sensitive functionality
* The lack of authentication could lead to unauthorized access to protected resources
* The API design violates basic security principles like "authentication before authorization"

### Recommended Fixes:
* Implement proper authentication (e.g., JWT, session tokens)
* Add authorization checks for sensitive endpoints
* Implement rate limiting
* Remove exposure of secret commands in the options endpoint
* Add API versioning and proper input validation

This type of vulnerability would typically be classified as CWE-306 (Missing Authentication for Critical Function) and could be considered OWASP API Security Top 10 2023 - API1:2023 Broken Object Level Authorization.

## Flag

```
HTB{D3v3l0p3r_t00l5_4r3_b35t__t0015_wh4t_d0_y0u_Th1nk??}
```
