from base64 import decode
from email import message
from bottle import get, post, route, response, run, request
from xml.dom import minidom
import redis
import json
import hashlib
import uuid
import xml.etree.ElementTree as ET

users = {
    "12345": {"id": "1", "email": "@a", "token": "12345"},
    "67890": {"id": "2", "email": "@b", "token": "67890"}    
}

#Use UUID1 or timestamps to compare

messages = {
    "1": [
        {"id": "1ee3356f-aadc-41d4-b75c-fd44896e9d13", "message": "Hello World 1", "access": "*", "Created_at": "insert time here"},
        {"id": "d20f9d1a-e345-4e70-a75a-5c3fd56f1f05", "message": "Hello World 2", "access": "*", "Created_at": "insert time here"},
        {"id": "a32b681b-e478-43bd-a501-8539d82e892a", "message": "Hello World 3", "access": "*", "Created_at": "insert time here"},
        {"id": "e4cbd305-ef9a-4c18-b6bc-ff55dd702e3d", "message": "Hello World 4", "access": "*", "Created_at": "insert time here"}
    ]
}


@get("/provider/<id>/from/<last_message_id:int>/limit/<limit:int>/token/<token>")
def _(id, last_message_id, limit, token):
    try:
        response.content_type = "applications/json"

        #Validation of limit
        if limit <= 0: raise Exception("Limit can't be 0 or below")

        #Validation of valid token
        if token not in users: raise Exception("Token is invalid!")

        #For Redis
        # keyList = r.keys()
        # redismessages = []
        
        # for rediskey in keyList[last_message_id:limit+last_message_id]:
        #     print(str(rediskey.decode("utf-8")))
        #     decodekey = str(rediskey.decode("utf-8"))
        #     redismessages.append({decodekey: r.get(decodekey).decode("utf-8")})
        
        # return json.dumps(redismessages)

        if last_message_id == 0:
            return json.dumps(messages[id][:limit])


        for dic in messages[id]:
            print(dic["id"])
            if dic["id"] == last_message_id:
                return json.dumps(messages[id][:limit])
                    
    except Exception as ex:
        response.status = 400 #Bad Request
        return str(ex)

@route("/createMessage/format/<format>/token/<token>", method="POST")
def _(format, token):
    try:
        #Hardcoded token received after MitID validation
        if(token != "12345"): raise Exception("Token is invalid!")
        if(format == "JSON"):

        #Hash both the ID and the topic into the key, likewise hash message
            id = uuid.uuid1()
            topic = request.json.get("topic")
            joinedkeyvalue = str(id) + ":" + str(topic)
            
            hash_object = hashlib.sha1(str.encode(joinedkeyvalue))
            hex_dig = hash_object.hexdigest()
            print(hex_dig)
            
            jsonMessage = request.json.get("message")

            #Send to redis
            #r.mset({hex_dig: jsonMessage})
            
            print("Topic is " + str(r.get(hex_dig).decode("utf-8")))

        if(format == "XML"):
            xmlRequest = request.body.read().decode("utf-8")
            print(xmlRequest)
            xmldoc = minidom.parseString(xmlRequest)
            root = ET.parse(xmldoc).getroot
            print(root)

        if(format == "YAML"):
            print(request.body)
            print("YAML")

        if(format == "TSV"):
            print(request.body)
            print("TSV")


    except Exception as ex:
        response.status = 400
        return str(ex)

#Connect to redis
#r = redis.Redis()
run(host="127.0.0.1", port=3000, debug=True, reloader=True)

