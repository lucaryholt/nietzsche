from get_name import name
from get_last_name import last_name
from get_email import email
from get_phone import phone
from get_api_key import api_key

import requests

message = f"Hello {name} {last_name} - Email: {email}"
print(phone)
print(message)
payload = {"to_phone": phone, "message": message, "api_key": api_key}
r = requests.post(
    "https://fatsms.com/send-sms",
    data=payload
)
print(r.status_code)
print(r.text)