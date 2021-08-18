# On-the-fly Image Resizer
This service can resize nor convert images on-the-fly and is designed to work with AWS S3.

## Usage
The **endpoint** for manipulating images is **/resize**.  
The service parses a query string for params to work with; here some working examples:
```
?key=images/theshot-640x0.jpeg
```
Will resize to a width of 640 pixels and keep the aspect ratio. The resulting redirect URI (301) will be:  
*images/theshot-640x0.jpeg*

```
?key=blue_marble-720x720.jpg
```
Will resize the image to 720 by 720 in width and height. The resulting redirect URI (301) will be:  
*blue_marble-720x720.jpg*

```
?key=gopher-0x480-jpg.png
```
Will resize the image to a height of 480 and keep the aspect ratio for width. The resized image is then converted from png to jpg. 
The resulting redirect URI will be:  
*gopher-0x480.jpg*

## Possible Improvements
- Replace github.com/disintegration/imaging with github.com/h2non/bimg

## Prerequisites
You'll need a ready-to-use AWS account as well as a S3 bucket in place. 
Create a bucket policy to allow anonymous access. 
```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "PublicRead",
            "Effect": "Allow",
            "Principal": "*",
            "Action": [
                "s3:GetObject",
                "s3:GetObjectVersion"
            ],
            "Resource": "arn:aws:s3:::___BUCKET NAME___/*"
        }
    ]
}
```
You have to replace the ___BUCKET NAME___ placeholder with the your bucket name.

The S3 bucket must be configured to host a static website. For the index document just enter index.html.
In the website hosting you need to edit Redirection Rules like so:
```
[
    {
        "Condition": {
            "HttpErrorCodeReturnedEquals": "404"
        },
        "Redirect": {
            "HostName": "___HOSTNAME___",
            "HttpRedirectCode": "307",
            "Protocol": "https",
            "ReplaceKeyPrefixWith": "resize?key="
        }
    }
]
```
You have to replace the ___HOSTNAME___ placeholder with the host, where the Image Resizer is running on.
