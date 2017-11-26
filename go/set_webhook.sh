#!/bin/bash

# hi, use this command to set webhook

curl -F "url=[url]" -F "certificate=[cert_path]" https://api.telegram.org/bot[token]/setWebhook
