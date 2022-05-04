from calendar import c
from pydoc import render_doc
import time
import random
from email.mime import application
from bottle import run, static_file, request, get, route, post, view, default_app
import requests
import mysql.connector
from mysql.connector import Error
import jwt

#################################


@get("/")
@view("index")
def _():
    return

@route("/validate", method="POST")
def validate():
    token = request.json.get("jwt")
    try:
        decoded_jwt = jwt.decode(token, "secret", algorithms=["HS256"])
        print(decoded_jwt)
        cpr = decoded_jwt.get("cpr")
        iat = decoded_jwt.get("iat")
        exp = decoded_jwt.get("exp")
        code = random.randint(0, 9999)
        print("CPR: " + cpr)
        print("Code: " + str(code))

        try:
            #mysql_insert_query = """SELECT * FROM verifications WHERE code = %s"""
            #record = (code,)
            mysql_insert_query = """INSERT INTO verifications (cpr, iat, exp, code, jwt)
                                    VALUES (%s, %s, %s, %s, %s)"""
            record = (cpr, iat, exp, code, str(token))
            
            cursor = connection.cursor()
            cursor.execute(mysql_insert_query, record)
            connection.commit()

            print("Record inserted")
            return "<p>Happy Times! the token is valid </p>"


        except mysql.connector.Error as error:
            print(error)

    except jwt.DecodeError:
        print("Invalid Token...")
        

@route("/createValidation", method="POST")
def validate():
    cpr = request.json.get("cpr")
    iat = request.json.get("iat")
    exp = request.json.get("exp")
    code = random.randint(0, 9999)
    print(code)

    try:
        mysql_insert_query = """INSERT INTO verifications (cpr, iat, exp, code)
                                VALUES (%s, %s, %s, %s)"""
        record = (cpr, iat, exp, code)
        cursor = connection.cursor()
        cursor.execute(mysql_insert_query, record)
        connection.commit()
        print("Record Inserted")

        sendSMS(22323386, code)

    except mysql.connector.Error as error:
        print("Error: ")
        print(error)


def sendSMS(phone, code):
    print("Attempting to send a text...")
    API_ENDPOINT = "https://fatsms.com/send-sms"
    API_KEY = "189c70a8-d6de-48c0-8553-97f03d8eaf89"
    body = {"to_phone": phone, "message": code, "api_key": API_KEY}
    response = requests.post(url=API_ENDPOINT, data=body)
    print(response)
    print(response.content)


#################################
try:
    #Server - Prod
    import production
    application = default_app()
except:
    #Localhost - Dev
    try:
        connection = mysql.connector.connect(
            host="localhost",
            database="systemsintegration",
            user="root",
            password="CsaE2PTRLxle#cYThl08br^6")
        if connection.is_connected():
            db_Info = connection.get_server_info()
            print("-:- Connected to MySQL Server version ", db_Info)
    except Error as e:
        print("-:- Error while connecting to MySQL", e)

    run(host="127.0.0.1", port=3333, debug=True, reloader=True, server="paste")
