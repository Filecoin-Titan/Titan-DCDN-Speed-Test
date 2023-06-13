# Titan-Dcdn-Speed-Test

`titan-dcdn-speed-test` is a command-line interface tool for interacting with the titan network for testing download speed.
Use the API provided by [titan-sdk-go](https://github.com/Filecoin-Titan/titan-sdk-go).

## Usage
```
NAME:
   titan-dcdn-speed-test - titan network toolset

USAGE:
   titan-dcdn-speed-test [global options] command [command options] [arguments...]

COMMANDS:
   download, d  Get file from Titan network
   test, t      Test the bandwidth of nodes with specified resources
   run          Start Titan speed test server
   help, h      Shows a list of commands or help for one command

```

To test the speed of the single L2 nodes, enter the following command in your terminal:

```shell
  titan-dcdn-speed-test test -c abcd
```
This will initiate a speed test for each node that has access to the resource. The results of the test will be displayed in your terminal.

To test the speed of fetching data concurrently from multiple nodes by using SDK, enter the following command in your terminal, specify the CID with the -c flag, 
and the output file location with the -o flag, -v to show more download details:

```
    titan-dcdn-speed-test download -c abcd
```

## Testing
The following tests will utilize the same resource, `QmTf6bbSW3rSaLuYPXxhPzAawe4iWwH1goKoXy5GxgZAwf`.

We will first test the speed of each L2 node.

```shell

./titan-dcdn-speed-test test -c QmTf6bbSW3rSaLuYPXxhPzAawe4iWwH1goKoXy5GxgZAwf
Start testing node speed ...
Direct Connection Nodes Time to First Byte(TTFB): 82ms
Nat Traversal (Punching) Nodes Time to First Byte(TTFB): 4402ms
+------------------------------------+-----------------------+------------+
|                NODE                |          IP           |   SPEED    |
+------------------------------------+-----------------------+------------+
| e_50acba1097d94387bd2af5b277ff4c17 | 121.57.52.48:33822    | 4.291MiB/s |
| e_1b77adf41afb4a2482f2e8c1f2c458ac | 123.134.50.22:1234    | 2.17MiB/s  |
| e_31cd3756a1594ba79c8de22baea8ab3c | 118.73.40.255:13354   | 3.531MiB/s |
| e_4f64056205c148e086bab5e0296af30f | 27.209.30.7:1234      | 2.704MiB/s |
| e_f6ba1f8fd4f64b8d98e28614429813ab | 124.134.231.250:22858 | 2.291MiB/s |
| e_31d56cb66d104a3fb13cc97d290a08ce | 39.79.173.188:1234    | 3.568MiB/s |
| e_953050e65db04eb39bbe43263d6fd8c2 | 112.245.153.240:1234  | 3.973MiB/s |
| e_4c4649f7a7f646b19ba2150eb79752ec | 27.214.89.39:26683    | 3.709MiB/s |
| e_6e9ed74d990b481daf09b799942d1380 | 110.182.42.155:1234   | 2.727MiB/s |
| e_e47c56cb127a4e41ba066bdcb45fb31a | 39.79.173.75:1234     | 3.68MiB/s  |
| e_819c3b30f0e84e8aaf94c054408a58c1 | 112.245.187.217:1234  | 3.697MiB/s |
| e_e6a953e88af4486c836d20d3012252bb | 119.186.183.211:47153 | 2.954MiB/s |
| e_eb6f338d3e28437684f831bbe3db1c56 | 39.79.228.208:1234    | 3.252MiB/s |
| e_6c17e889d82049de9c53b52f00127e0d | 39.79.226.3:1234      | 1.439MiB/s |
| e_829620625abf43cfb8d265e64ead5273 | 119.187.155.173:1234  | 4.483MiB/s |
| e_4bfdb74b46cd408c961b389f51231d36 | 110.182.212.95:1234   | 998.1KiB/s |
| e_56885f67882d420f804b94d2f55184d5 | 60.214.245.199:1234   | 2.95MiB/s  |
| e_a1b05a6568654f11974bfb1c983e850b | 112.243.109.73:1234   | 3.283MiB/s |
| e_e52f92ff8c85461aaf2d9a89915920cb | 60.214.227.40:1234    | 2.94MiB/s  |
| e_ee1deff728454c0295c7ae169d2c6ca0 | 112.21.18.71:4821     | 2.757MiB/s |
| e_777dc64f219e4fdeaf6669d733f4a717 | 223.149.14.15:50629   | 3.724MiB/s |
| e_c03790f69c8f48fe9535dd6eb817733f | 111.16.253.242:1553   | 3.313MiB/s |
| e_83434fbb56e94f9097a0ed1125957e50 | 60.214.219.71:1234    | 3.515MiB/s |
| e_c120bcaf88cd47c398a1467a998e27d6 | 39.79.169.126:1234    | 4.149MiB/s |
| e_d35557b49b034922904ad67428555061 | 223.96.113.122:11804  | 1.003MiB/s |
| e_3bb2cdf041f348bbb9d6f5ce63805b1a | 111.227.253.92:25797  | 4.356MiB/s |
+------------------------------------+-----------------------+------------+



```

Next, we utilize the SDK to download the same resource.

```shell
./titan-dcdn-speed-test download -v -c QmTf6bbSW3rSaLuYPXxhPzAawe4iWwH1goKoXy5GxgZAwf
+------------------------------------+-----------------------+------------+-------+----------+
|               NODEID               |        ADDRESS        |   SPEED    | COUNT | DATASIZE |
+------------------------------------+-----------------------+------------+-------+----------+
| e_56885f67882d420f804b94d2f55184d5 | 60.214.245.199:1234   | 470KiB/s   |     2 | 4MiB     |
| e_1b77adf41afb4a2482f2e8c1f2c458ac | 123.134.50.22:1234    | 2.379MiB/s |     8 | 16MiB    |
| e_31cd3756a1594ba79c8de22baea8ab3c | 118.73.40.255:13354   | 1.095MiB/s |     4 | 8MiB     |
| e_eb6f338d3e28437684f831bbe3db1c56 | 39.79.228.208:1234    | 2.18MiB/s  |     8 | 16MiB    |
| e_e47c56cb127a4e41ba066bdcb45fb31a | 39.79.173.75:1234     | 853KiB/s   |     3 | 6MiB     |
| e_f6ba1f8fd4f64b8d98e28614429813ab | 124.134.231.250:22858 | 451.2KiB/s |     1 | 2MiB     |
| e_c03790f69c8f48fe9535dd6eb817733f | 111.16.253.242:1553   | 1.43MiB/s  |     5 | 10MiB    |
| e_6e9ed74d990b481daf09b799942d1380 | 110.182.42.155:1234   | 1.141MiB/s |     4 | 8MiB     |
| e_953050e65db04eb39bbe43263d6fd8c2 | 112.245.153.240:1234  | 767.7KiB/s |     3 | 6MiB     |
| e_819c3b30f0e84e8aaf94c054408a58c1 | 112.245.187.217:1234  | 544.1KiB/s |     1 | 2MiB     |
| e_4f64056205c148e086bab5e0296af30f | 27.209.30.7:1234      | 547.6KiB/s |     1 | 2MiB     |
| e_4bfdb74b46cd408c961b389f51231d36 | 110.182.215.250:1234  | 508.6KiB/s |     1 | 2MiB     |
| e_e6a953e88af4486c836d20d3012252bb | 119.186.183.211:47153 | 494.4KiB/s |     1 | 2MiB     |
| e_a1b05a6568654f11974bfb1c983e850b | 112.243.109.73:1234   | 2.079MiB/s |     7 | 14MiB    |
| e_c120bcaf88cd47c398a1467a998e27d6 | 39.79.169.126:1234    | 1.322MiB/s |     5 | 9.93MiB  |
| e_83434fbb56e94f9097a0ed1125957e50 | 60.214.219.71:1234    | 1.06MiB/s  |     4 | 8MiB     |
| e_5d210f13d98d4b40a3a695d25a542d2a | 39.79.227.111:1234    | 1.141MiB/s |     4 | 8MiB     |
| e_3bb2cdf041f348bbb9d6f5ce63805b1a | 111.227.253.92:25797  | 873.5KiB/s |     2 | 4MiB     |
| e_4c4649f7a7f646b19ba2150eb79752ec | 27.214.89.39:26683    | 811.4KiB/s |     1 | 2MiB     |
| e_e52f92ff8c85461aaf2d9a89915920cb | 60.214.227.40:1234    | 768.3KiB/s |     3 | 6MiB     |
| e_31d56cb66d104a3fb13cc97d290a08ce | 39.79.173.188:1234    | 486.1KiB/s |     1 | 2MiB     |
| e_50acba1097d94387bd2af5b277ff4c17 | 121.57.52.48:33822    | 1.413MiB/s |     6 | 10MiB    |
| e_777dc64f219e4fdeaf6669d733f4a717 | 223.149.154.131:63745 | 1.706MiB/s |     6 | 12MiB    |
| e_ee1deff728454c0295c7ae169d2c6ca0 | 112.21.18.71:4821     | 1.112MiB/s |     4 | 8MiB     |
| e_d35557b49b034922904ad67428555061 | 223.96.113.122:11804  | 1.177MiB/s |     5 | 10MiB    |
+------------------------------------+-----------------------+------------+-------+----------+
 177.93 MiB / 177.93 MiB [==============================================================================================================================] 100.00% 18.33 MiB/s 9s

```
The above test shows that the SDK retrieve data from multiple nodes, and significantly improves download speed.
The final download speed is subject to the local network bandwidth.