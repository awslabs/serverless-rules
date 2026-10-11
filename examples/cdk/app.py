#!/usr/bin/env python3

from aws_cdk import App

from cdk.cdk_stack import CdkStack


app = App()
CdkStack(app, "CdkStack")

app.synth()
