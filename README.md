# commitgrabber

While enumerating a web application, it can be hard to find what is the framework version in use.
This can happen when there is no information in headers, cookies, config files, etc.
An alternative approach is to fingerprint the web application using files that are present in the framework repository (i.e some frameworks have `.js` files that are available in the repository and in the web application). Knowing this it is possible to extract the md5sum of that file and find (in the commit logs), what is the commit hash of that specific file.

commitgrabber aims to facilitate this process by searching in the commit logs what version corresponds a specific md5sum.

## PoC

- Curl a `.js` file from the web app get the md5sum

``` bash
$ curl https://something/some.js -o some.js ; md5sum some.js
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 22695  100 22695    0     0  81417      0 --:--:-- --:--:-- --:--:-- 81344
1f7ca2086ba90f76bf684edc65fc89a1  some.js
```

- Navigate to the web app repo and curl the `.js` file that should be in the latest version
  - Note that if the md5sum is the same, that means the web application is using the latest version

``` bash
$ curl https://raw.githubusercontent.com/Something/master/public/js/some.js -o somefromrepo.js && md5sum somefromrepo.js 
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 23039  100 23039    0     0   339k      0 --:--:-- --:--:-- --:--:--  340k
e0892b7039a41978f37e2c2cc0a5180b  somefromrepo.js
```

- Execute commitgrabber to find the commit link in order to get the version in use
  - url is the .git repo in HTTPS
  - hash is the md5sum from the first curl
  - fn is the file name that we are searching for

``` bash
./commitgrabber -url "https://github.com/Something/s.git" -hash "1f7ca2086ba90f76bf684edc65fc89a1" -fn "some.js"

Commit link: https://github.com/Something/some/commit/b061bb13872806e6a3fe3a5227961e18e34c2839

commit b061bb13872806e6a3fe3a5227961e18e34c2839
Author: NoOne <no.one@gmail.com>
Date:   Tue Jul 25 16:00:55 2021 +0100

    added stuff to javascript file
```

- Note that commitgrabber does not write to disk so it might take a while depending on the size of the repository

- Alernative you can use the following one liner in bash

``` bash
for i in $(git log some.js | grep ^commit | awk '{print $2}'); do git checkout $i -- some.js; echo -n "$i "; md5sum some.js; done
```

## TODO

- [x] Add README
- [] Refactor code
- [] Add Tests
