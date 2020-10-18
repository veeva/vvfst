# Commands
## Login
This is one of the first commands and is a must to login to the vault.  This login command accepts required details to login.  The login information and response are cached locally in user's home directory with filename, $HOME/.vvfst.yaml.

#### Usage
```
vvfst login --help
Login with username and password for the vault.
For example:
  login --domain_name myvalut.veevavault.com --login myuser@mydomain.com
  login -d myvault.veevavault.com -a v20.1 -u myuser@mydomain.com

Usage:
  vvfst login [flags]

Flags:
  -a, --api_version string   API Version 
  -d, --domain_name string   Vault domain name 
  -h, --help                 help for login
  -u, --username string      Vault username 

Global Flags:
  -x, --debug   Enable debug
```  

* These information are locally cached, hence the flags are not required for subsequent login.  
* These information are used to auto login if REST API session is expired.  The cli provides infinite session experience.
* The login -h will display all cached information, no need to hunt for the file.
* After successful login, it also provides duration that it took to execute the REST api.

#### Examples:
```
vvfst login -a v20.2 -d mylogin.vaultdev.com -u myuser@mydomain.com
10:27AM INFO  [Duration: 2.294 seconds] Login successful.

vvfst login10:28AM INFO [Duration: 2.132 seconds] Login successful.
```
Note: The login information and response cached under the $HOME/.vvfst.yaml and you may view them by cat $HOME/.vvfst.yaml


## Logout
The user can logout after the session is done.  It also provides an option to clean up (purge) all information cached locally including the file $HOME/vvfst.yaml.

#### Usage
```
vvfst logout --help
Logout and delete all cached session data.

Usage:
  vvfst logout [flags]

Flags:
  -c, --clear   Clear all configuration data
  -h, --help    help for logout

Global Flags:
  -x, --debug   Enable debug
```

#### Examples:
```
vvfst logout
10:37AM INFO  Logout successful.

vvfst logout -c
Clearing configuration..
```

## List
Listing a directory is one of the basic functionality and it helps to visualize what is stored in the file staging area.  By default, it lists all files in the user's home directory, user can specify any directory as well.

#### Usage
```
vvfst ls --help
List the content in the directory.  The listing is a flat list when including sub directories.

Usage:
  vvfst ls <remote-file/folder> [flags]

Flags:
  -h, --help        help for ls
  -l, --limit int   Limit number of items (default 100)
  -r, --recursive   Enable recursive mode to list all sub directories

Global Flags:
  -x, --debug   Enable debug
```

#### Examples:
````
vvfst ls
11:08AM INFO  [Duration: 0.810 seconds] ls completed.
listing: /
kind    path                                                size
====================================================================
folder  /Inbox/                                             0 B
folder  /aws/                                               0 B
file    /.DS_Store                                          10.2 kB
file    /a.txt                                              69.7 MB

## To list specific diretory
vvfst ls /aws
11:14AM INFO  [Duration: 1.019 seconds] ls completed.
listing: /aws
kind    path                                                size
====================================================================
folder  /aws/aws-cli/                                       0 B
file    /aws/.DS_Store                                      10.2 kB


## To list content recursively
vvfst ls -r /
11:16AM INFO  [Duration: 1.037 seconds] ls completed.
listing: /
kind    path                                                size
====================================================================
file    /.DS_Store                                          10.2 kB
file    /a.txt                                              69.7 MB
file    /aws/.DS_Store                                      10.2 kB
file    /aws/aws-cli/.DS_Store                              8.2 kB
file    /aws/aws-cli/README.md                              2.3 kB
file    /aws/aws-cli/bash-linux/.DS_Store                   8.2 kB
file    /aws/aws-cli/bash-linux/README.md                   2.4 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  4.7 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  5.6 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  10.2 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  7.7 kB
file    /aws/aws-cli/bash-linux/s3/.DS_Store                8.2 kB
file    /aws/aws-cli/bash-linux/s3/bucket-lifecycle-oper..  3.1 kB
file    /aws/aws-cli/bash-linux/s3/bucket-lifecycle-oper..  5.5 kB


## To list content recursively and with limtit
vvfst ls / -r -l 5
11:17AM INFO  [Duration: 1.144 seconds] ls completed.
listing: /
kind    path                                                size
====================================================================
file    /.DS_Store                                          10.2 kB
file    /a.txt                                              69.7 MB
file    /aws/.DS_Store                                      10.2 kB
file    /aws/aws-cli/.DS_Store                              8.2 kB
file    /aws/aws-cli/README.md                              2.3 kB
11:17AM INFO  Press Ctl+C to stop or press space-bar for next page
11:17AM INFO  [Duration: 0.420 seconds] ls completed.
listing: /
kind    path                                                size
====================================================================
file    /aws/aws-cli/bash-linux/.DS_Store                   8.2 kB
file    /aws/aws-cli/bash-linux/README.md                   2.4 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  4.7 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  5.6 kB
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  10.2 kB
11:17AM INFO  Press Ctl+C to stop or press space-bar for next page
11:17AM INFO  [Duration: 0.455 seconds] ls completed.
listing: /
kind    path                                                size
====================================================================
file    /aws/aws-cli/bash-linux/ec2/change-ec2-instance-..  7.7 kB
file    /aws/aws-cli/bash-linux/s3/.DS_Store                8.2 kB
file    /aws/aws-cli/bash-linux/s3/bucket-lifecycle-oper..  3.1 kB
file    /aws/aws-cli/bash-linux/s3/bucket-lifecycle-oper..  5.5 kB

## To download csv reports
vvfst ls -c -r
10:54PM INFO  Current job status: RUNNING
10:54PM INFO  / list export as csv
10:54PM INFO  Downloading reports 123611.csv

cat 123611.csv
"kind","path","name","size","modified_date"
"file","/b","b",69650794,"2020-10-20T00:51:48.000Z"
````

## Create Directory
A best way to organize files is to create a directory and keep the content inside a directory.  The cli allows user to create a directory.

#### Usage
```
vvfst mkdir --help
Create remote directory if not exists

Usage:
  vvfst mkdir <remote-folder> [flags]

Flags:
  -h, --help        help for mkdir
  -o, --overwrite   Enable overwrite to overwrite existing directory

Global Flags:
  -x, --debug   Enable debug
```    


## Upload 
This is a most important stuff, the user wants to upload content so they can load them into the vault.  It helps to upload a single small file, any large file or even any folder with any number of content including large file.  Upload also helps to upload concurrently using multiple threads.  It also handles the large file inside the folder automatically.

#### Usage
```
vvfst upload --help
Uploading a single file or all files from a folder

Usage:
  vvfst upload <local-file/folder> <remote-file/folder> [flags]

Flags:
  -h, --help              help for upload
  -o, --overwrite         Enable overwrite to overwrite if file/folder exists
  -t, --threadCount int   Number of concurrent thread to upload (default 1)

Global Flags:
  -x, --debug   Enable debug
```

one file uses only one thread, multiple thread is not going to increase speed for a single file.


#### Examples
````
## Uploading single file
vvfst upload ~/tmp/demo3/aws/aws-cli/README.md /Readme.txt
11:39AM INFO  [Duration: 1.037 seconds] uploaded file: /Readme.txt

## Uploading a single large file
vvfst upload ~/tmp/demo3/consoleText.txt /example.txt
11:38AM INFO  [Duration: 0.466 seconds] upload session created for file: /example.txt
11:38AM INFO  [/example.txt] Uploaded part: 1 of 14, size: 5.2 MB, partContentMD5: be8b20436c596cb309c32cbd2afb8e56
11:38AM INFO  [/example.txt] Uploaded part: 2 of 14, size: 5.2 MB, partContentMD5: c6e94d4eac86cb99de4f39859d05ef46
....
11:39AM INFO  [/example.txt] Uploaded part: 14 of 14, size: 1.5 MB, partContentMD5: 1edb256cdb25bba96e54f75c01ac3b5f
11:39AM INFO  [Duration: 0.230 seconds] upload session completed for file: /example.txt, waiting for job completion
11:39AM INFO  Current job status: QUEUED
11:39AM INFO  /example.txt file upload sucessfully

## Uploading a directory 
vvfst upload ~/tmp/demo3 /demo3
11:41AM INFO  [Duration: 1.165 seconds] uploaded file: /demo3/.DS_Store
11:41AM INFO  [Duration: 0.565 seconds] uploaded file: /demo3/aws/.DS_Store
11:41AM INFO  [Duration: 0.807 seconds] uploaded file: /demo3/aws/aws-cli/.DS_Store
11:41AM INFO  [Duration: 0.507 seconds] uploaded file: /demo3/aws/aws-cli/README.md
11:41AM INFO  [Duration: 0.529 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/.DS_Store
11:41AM INFO  [Duration: 0.773 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/README.md
11:41AM INFO  [Duration: 0.493 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/README.md
11:41AM INFO  [Duration: 0.499 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/awsdocs_general.sh
11:41AM INFO  [Duration: 0.618 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/change_ec2_instance_type.sh

## Upload a directory with multiple thread and overwrite 
vvfst upload ~/tmp/demo3 /demo3 -t 10 -o
11:42AM INFO  [Duration: 1.159 seconds] uploaded file: /demo3/.DS_Store
11:42AM INFO  [Duration: 1.138 seconds] uploaded file: /demo3/aws/.DS_Store
11:42AM INFO  [Duration: 0.999 seconds] uploaded file: /demo3/aws/aws-cli/README.md
11:42AM INFO  [Duration: 1.031 seconds] uploaded file: /demo3/aws/aws-cli/.DS_Store
11:42AM INFO  [Duration: 0.823 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/.DS_Store
11:42AM INFO  [Duration: 1.001 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/README.md
11:42AM INFO  [Duration: 0.755 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/README.md
11:42AM INFO  [Duration: 0.785 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/change_ec2_instance_type.sh
11:42AM INFO  [Duration: 0.826 seconds] uploaded file: /demo3/aws/aws-cli/bash-linux/ec2/change-ec2-instance-type/test_change_ec2_instance_type.sh
......
11:42AM INFO  [Duration: 0.505 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_amazonaws_aws_java_sdk_core_1_11_738.xml
11:42AM INFO  [Duration: 0.506 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_amazonaws_jmespath_java_1_11_738.xml
11:42AM INFO  [Duration: 0.509 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_fasterxml_jackson_core_jackson_databind_2_6_7_3.xml
11:42AM INFO  [Duration: 0.511 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_amazonaws_aws_java_sdk_kms_1_11_738.xml
11:42AM INFO  [Duration: 0.527 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_fasterxml_jackson_core_jackson_core_2_6_7.xml
11:42AM INFO  [Duration: 0.569 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__com_fasterxml_jackson_dataformat_jackson_dataformat_cbor_2_6_7.xml
11:42AM INFO  [Duration: 0.483 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__commons_codec_commons_codec_1_14.xml
11:42AM INFO  [Duration: 0.492 seconds] uploaded file: /demo3/aws/s3/.idea/libraries/Maven__commons_io_commons_io_2_6.xml
.....
11:43AM INFO  [Duration: 0.500 seconds] uploaded file: /demo3/aws/s3/testfile1.txt
11:43AM INFO  [Duration: 0.339 seconds] upload session created for file: /demo3/consoleText.txt
11:43AM INFO  [Duration: 0.636 seconds] uploaded file: /demo3/docs-2020-10/FileStagingAPI-Page-3.png
11:43AM INFO  [/demo3/consoleText.txt] Uploaded part: 1 of 14, size: 5.2 MB, partContentMD5: be8b20436c596cb309c32cbd2afb8e56
11:43AM INFO  [/demo3/consoleText.txt] Uploaded part: 2 of 14, size: 5.2 MB, partContentMD5: c6e94d4eac86cb99de4f39859d05ef46
11:43AM INFO  [/demo3/consoleText.txt] Uploaded part: 3 of 14, size: 5.2 MB, partContentMD5: 175933b42e5e02f92cbdfed5f86ba53a
.....
11:43AM INFO  [/demo3/consoleText.txt] Uploaded part: 14 of 14, size: 1.5 MB, partContentMD5: 1edb256cdb25bba96e54f75c01ac3b5f
11:43AM INFO  [Duration: 0.206 seconds] upload session completed for file: /demo3/consoleText.txt, waiting for job completion
11:43AM INFO  Current job status: RUNNING
11:43AM INFO  /demo3/consoleText.txt file upload successfully
````

## Download
The cli tool allows user to download a single file or directory.  It also helps to download concurrently.  By default, it downloads only items from the source folder and recursive mode allows to download entire folder.  It also has a nice download progress bar.

#### Usage
````
vvfst download --help
Download folder/files from remote staging folder to current folder

Usage:
  vvfst download <remote-file/folder> <local-file/folder> [flags]

Flags:
  -h, --help              help for download
  -r, --recursive         Enable recursive mode to download all sub directories
  -t, --threadCount int   Number of concurrent thread to download (default 1)

Global Flags:
  -x, --debug   Enable debug
````

Note: Only last progressbar get updated since console output does not have a great way to update multiple lines or past line at the same time.



#### Examples
````
## Donwloading from home directory
vvfst download / /tmp/a
downloading .DS_Store  100% >==================================================================================================================| (10244/10244, 51772412 it/s) [0s:0s]
downloading Readme.txt  100% >===================================================================================================================| (2303/2303, 10483144 it/s) [0s:0s]
downloading a.txt   45% >==================================================                                                              | (31916032/69650794, 24627267 it/s) [1s:1s]

## Downloading recursively with multiple thread 
vvfst download / /tmp/a -t 10 -r
downloading .DS_Store  100% >==================================================================================================================| (10244/10244, 46555807 it/s) [0s:0s]
downloading README.md  100% >====================================================================================================================| (4669/4669, 21915146 it/s) [0s:0s]
````

## Move
Move a file from one folder to another folder.  Even move files from one name to another name, oops! it is called rename.   It applies to directory as well.  This command lets you move a directory to another directory or rename.  This move or rename is invoking the REST API which creates an asynchronous job for the actual move, the commands waits for up to 1 minute for the job to complete and the job will run after 1 minute even if the command quits.

#### Usage
```
vvfst mv --help
Move files/folder from one location another location within remote location or rename file/folder in the remote directory

Usage:
  vvfst mv <src-remote-file/folder> <dest-remote-file/folder> [flags]

Flags:
  -h, --help        help for mv
  -o, --overwrite   Enable overwrite to overwrite a file or merge existing folders

Global Flags:
  -x, --debug   Enable debug
```
Examples
```
## Move a file
vvfst mv /Readme.txt /Hello.txt
12:05PM INFO  [Duration: 1.192 seconds] mv submitted successfully, waiting for job completion
12:05PM INFO  Current job status: RUNNING
12:05PM INFO  /Readme.txt moved to /Hello.txt successfully

## Move a directory
vvfst mv /demo3/ /demo3a
12:12PM INFO  [Duration: 1.196 seconds] mv submitted successfully, waiting for job completion
12:12PM INFO  Current job status: RUNNING
12:12PM INFO  Current job status: RUNNING
12:12PM INFO  Current job status: RUNNING
12:12PM INFO  Current job status: RUNNING
12:12PM INFO  Current job status: RUNNING
12:13PM INFO  Current job status: RUNNING
Error: Job not completed within 60 seconds
```

## Delete
If incorrect files are uploaded, then we should give the user an option to delete.  This tool provides delete command.  It deletes a single file or folder recursively.  It also executes REST API which creates asynchronous job, the command will wait up to 1 minute for job completion.

#### Usage
```
vvfst rm --help
Remove files/folder remote location

Usage:
  vvfst rm <src-remote-file/folder> <dest-remote-file/folder> [flags]

Flags:
  -h, --help        help for rm
  -r, --recursive   Enable recursive mode to delete all sub directory contents

Global Flags:
  -x, --debug   Enable debug
```  
#### Examples
```
## Delete a file
vvfst rm /a.txt
12:28PM INFO  [Duration: 0.985 seconds] rm submitted successfully, waiting for job completion
12:28PM INFO  Current job status: RUNNING
12:28PM INFO  /a.txt removed successfully

## Delete a folder recursively
vvfst rm /aws/ -r
11:40AM INFO  [Duration: 0.929 seconds] rm submitted successfully, waiting for job completion
11:40AM INFO  Current job status: RUNNING
11:40AM INFO  /aws/ removed successfully
```

## List Multipart upload session (Resumeable upload session)
Listing a multipart upload session is a must since this file is not visible to the user until it finishes the session.  The list details contains useful information such as upload parts, size, expiration , etc.

#### Usage
```
vvfst mls --help
List of multipart upload sessions initiated by that user which is not expired

Usage:
  vvfst mls <remote-file/folder> [flags]

Flags:
  -h, --help   help for mls

Global Flags:
  -x, --debug   Enable debug
```

#### Examples
```
## No session available
vvfst mls
12:31PM INFO  [Duration: 0.912 seconds] mListCmd completed.
listing upload sessions:
path                            size             up size          up parts   expiration
==============================================================================================

## when upload session is available
vvfst mls
12:33PM INFO  [Duration: 0.794 seconds] mListCmd completed.
listing upload sessions:
path                            size             up size          up parts   expiration
==============================================================================================
/a.txt                          69.7 MB          15.7 MB          3          2020-10-15 19:53
```

## Remove Multipart upload session (Resumable upload session)
The upload session is limited for each vault. The user must terminate any unused upload session so he/she can start new upload session.

#### Usage
```
vvfst mrm --help
Delete upload session started for the given file

Usage:
  vvfst mrm <remote-file> [flags]

Flags:
  -h, --help   help for mrm

Global Flags:
  -x, --debug   Enable debug
```
#### Examples
```
vvfst mrm /example.txt
11:38AM INFO  [Duration: 0.272 seconds] Deleted upload session for /example.txt
```

## Jobs 
To check active jobs created using the cli, can be tracked.  Once job completes then it will be removed from active job list.

#### Usage
```
vvfst jobs --help
Display list of jobs and validate job status.  If job is completed and it will be removed from the list.  It keep track jobs submitted via this cli from this computer.

Usage:
  vvfst jobs [flags]

Flags:
  -h, --help                help for jobs
  -t, --threadCount int     Number of concurrent thread to check job status (default 1)
  -T, --timoutSeconds int   How long job status to be checked (default 60)

Global Flags:
  -x, --debug   Enable debug
```

#### Examples
```
## When no jobs are available
vvfst jobs
10:57PM INFO  No active job(s) available

## When jobs are running
vvfst mv /b /c
10:57PM INFO  [Duration: 1.272 seconds] mv submitted successfully, waiting for job completion
10:57PM INFO  Current job status: RUNNING
^C

vvfst jobs
10:57PM INFO  Checking job status: 123711
10:57PM INFO  /b moved to /c successfully
10:57PM INFO  Job completed -  /b moved to /c successfully
```



