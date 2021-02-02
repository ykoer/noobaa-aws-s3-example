import os
import logging
import boto3
from botocore.exceptions import ClientError

ACCESS_KEY = os.environ['AWS_ACCESS_KEY_ID']
SECRET_KEY = os.environ['AWS_SECRET_ACCESS_KEY']
BUCKET_HOST = os.environ['BUCKET_HOST']
BUCKET_NAME = os.environ['BUCKET_NAME']

s3_client = boto3.client('s3',
        verify=False,
        endpoint_url = BUCKET_HOST,
        aws_access_key_id=ACCESS_KEY,
        aws_secret_access_key=SECRET_KEY)

def put_object(bucket_name, key, value):
    try:
        s3_client.put_object(Bucket=bucket_name, Key=key, Body=value)
    except ClientError as e:
        logging.error(e)
        return False
    return True

def get_object(bucket_name, key):
    try:
        result = s3_client.get_object(Bucket=bucket_name, Key=key)
    except ClientError as e:
        logging.error(e)
        return None
    return result


def main():
    put_object(BUCKET_NAME, 'foo', 'bar')
    print(get_object(BUCKET_NAME, 'foo')['Body'].read())

if __name__ == '__main__':
    main()