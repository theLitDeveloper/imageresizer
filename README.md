# On-the-fly Image Resizer
This service can resize, crop nor convert images on-the-fly and is designed to work with AWS S3. The implementation was inspired by the idea from this blog [post on AWS](https://aws.amazon.com/de/blogs/compute/resize-images-on-the-fly-with-amazon-s3-aws-lambda-and-amazon-api-gateway/).

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
You have to replace the ```___BUCKET NAME___``` placeholder with the your bucket name.

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
            "ReplaceKeyPrefixWith": "do?ref="
        }
    }
]
```
You have to replace the ```___HOSTNAME___``` placeholder with the host, where the Image Resizer is running on.

### CloudFront

Next, you have to set up a CloudFront distribution.

Go to your newly created bucket, select Properties tab and scroll down to “Hosting a static web site”. Copy the endpoint URL without the scheme, e.g. “simplys3test.s3-website.eu-central-1.amazonaws.com“.

Open CloudFront in your AWS console. Paste the copied bucket endpoint URL, without scheme, in the “Origin domain” field. Then scroll down to “Viewer” and select “https only” for the Viewer protocol policy. Now click on “Create distribution”.

## Deployment / Development
Set environment variables:  
- AWS_ACCESS_KEY_ID: your AWS access key id
- AWS_SECRET_ACCESS_KEY: your AWS secret access key  
- REDIRECT_HOST: the domain from your CloudFront distribution
- AWS_BUCKET: your bucket name
- AWS_REGION: your region, e.g. eu-central-1

The default port for this service is 4321. You can easily adjust this by providing a PORT environment variable, e.g. PORT=8080

## Usage
The **endpoint** for manipulating images is **/do**. This endpoint expects a query string parameter ```ref``` with a URL string and contained parameters for changing the image.  

The endpoint for manipulating images is /do by default. The service parses the query string provided by your AWS S3 bucket. So a valid request look like  
```/do?ref=<params>/<resource path>/<resource name>```

Available parameters for image manipulation:  

**w_** (width) *required*  
Example: w_500  

**h_** (height) *required*  
Example: h_924  

**crop** *optional*  

**gray** *optional*  

Parameters can be specified in arbitrary order, only once each and separated by commas.
Example: ```h_500,w_500,gray/images/test.jpg```  
### Resizing
Resizing is the default action. If a 0 is passed for the width or the height, the image is resized in height or width while the aspect ratio is preserved.

### Cropping
Cropping works similarly to resizing, but instead of resizing the image, it is centered and cropped to the specified width and length.

### Convert to Grayscale
Just add ```gray``` to parameters to discard any color in your image and convert it to grayscale.

### Auto-conversion from PNG to JPG
This is done automatically when a PNG file is submitted to this service.