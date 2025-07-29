from classFile_finall import *
from hazelcast.predicate import ilike, sql, and_, between, equal, greater, less, in_, less_or_equal,greater_or_equal, or_
from hazelcast.serialization.api import Portable
from flask import Flask, jsonify
from flask_restful import Resource, Api, reqparse
from hazelcast.core import HazelcastJsonValue
from datetime import datetime,date
import hazelcast
import logging
import json
import time
import maps
from settings import db_settings


app = Flask(__name__)
api = Api(app)

logging.basicConfig(level=logging.INFO, filename='ReadFromPostgres.log', filemode='w', format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', datefmt='%Y-%m-%d %I:%M:%S %p')

logging.info("settings {}".format(db_settings))
import psycopg2
conn = psycopg2.connect(
      user = db_settings['user'],
      port= db_settings['port'],
      password = db_settings['password'],
      host = db_settings['host'],
      database = db_settings['database']
)
cur = conn.cursor()


#get Supplements_Hotel from db

sql = """SELECT id,period_id,room_id,type_id,board_id,adults,contrat_id,price,operation FROM db.supplements where deleted = false;"""
cur.execute(sql)
items= cur.fetchall()
for item in items :
    maps.supplements_hotel_map.put(item[0], Supplements_Hotel(item[0],item[1],item[2],item[3],item[4],item[5],item[6],item[7],item[8]))
logging.info("data_set Start **********")

#get Supplements_Pricing from db

sql = """SELECT opc.id, opc.contract_id, opc.options  FROM db.options_contract as opc join db.contract as ct on ct.id=opc.contract_id and ct.deleted=false"""
cur.execute(sql)
items= cur.fetchall()
for item in items :
    maps.supplement_princing_map.put(item[0], Supplement_princing(item[0],item[1],item[2]))


# #get client from db
#               0   1       2   3           4
# sql = """SELECT id,name, ordre, marge, marge_operation from db.client"""
# cur.execute(sql)
# clients= cur.fetchall()
# for clientt in clients :
#     maps.client_map.put(clientt[0], Client(clientt[0],clientt[1],clientt[2],clientt[3],clientt[4]))


#get group_admin from db
#               0   1       2          3        4       5        6
sql = """SELECT id,name,admin_max,hotel_max,agence_max,active,currency from db.group_admin"""
cur.execute(sql)
clients= cur.fetchall()
for clientt in clients :
    maps.group_admin_map.put(clientt[0], Group_Admin(clientt[0],clientt[1],clientt[2],clientt[3],clientt[4],clientt[5],clientt[6]))


#get admin from db
#               0   1    2     3        4       5
sql = """SELECT id,name,role,active,group_id,solde from db.admin where deleted = false"""
cur.execute(sql)
clients= cur.fetchall()
for clientt in clients :
    maps.admin_map.put(clientt[0], Admin(clientt[0],clientt[1],clientt[2],clientt[3],clientt[4],clientt[5]))


#get group_agence from db
#               0   1      2          3           4         5
sql = """SELECT id,name,marge, marge_operation,client_id,market from db.group_agence where deleted = false"""
cur.execute(sql)
group_agences= cur.fetchall()
for group_agence in group_agences :
    print("group",group_agence)
    maps.groupagence_map.put(group_agence[0], GroupAgence(group_agence[0],group_agence[1],group_agence[2],group_agence[3],group_agence[5],group_agence[4]))



#get users from db
#               0    1           2          3     4        5              6                  7              8     9
sql = """SELECT id,agence_id, active, group_id, solde, solde_rouge,  marge_b2b_operation, maxrequest, marge_b2b,currency from db.user where deleted = false"""
cur.execute(sql)
users= cur.fetchall()
for user in users :
    maps.user_map.put(user[0], User(user[0],user[1],user[2],user[3],user[4],user[5],user[6],user[7],user[8],user[9]))

#get agence from db
#               0     1      2            3         4             5       6           7       8              9        10
sql = """SELECT id, status, marge, marge_operation,solde, solde_rouge,currency,maxrequest,group_agence_id,marge_b2b,marge_xml from db.agence where deleted = false """
cur.execute(sql)
agences= cur.fetchall()
for agence in agences :
    maps.agence_map.put(agence[0], Agence(agence[0],agence[1],agence[2],agence[3],agence[4],agence[5],agence[6],agence[7],agence[8],agence[9],agence[10]))

#get countries from db
#                 0            1            2
sql = """SELECT countryid,countryname,countrycode from db.country """
cur.execute(sql)
countries= cur.fetchall()
for country in countries :
    maps.country_map.put(country[0], Country(country[0],country[1],country[2]))

#get cities from db
#                 0       1         2       3           4
sql = """SELECT cityid,citycode,cityname,countryid,countryname from db.city """
cur.execute(sql)
cities= cur.fetchall()
for city in cities :
    maps.city_map.put(city[0], City(city[0],city[1],city[2],city[3],city[4]))

#get hotels from db
#               0   1     2     3       4       5           6       7           8       9          10       11        12
sql = """SELECT id,name,city,country,stars,phone_number1,email1,description,latitude,longitude,address,hotelchain,client_id,short_description from db.hotel where deleted = false"""
cur.execute(sql)
hotels= cur.fetchall()
for hotel in hotels :
    stars = int(hotel[4]) if hotel[4] is not None else 0

    maps.hotels_map.put(hotel[0], Hotel(hotel[0],hotel[1],hotel[2],hotel[3],stars,hotel[10],hotel[6],hotel[8],hotel[9],hotel[5],hotel[7],hotel[11],hotel[12],hotel[13]))


#get hotelchains from db
#                0               1       2      3           4     5     6     7       8
sql = """SELECT hotelchainid,countryid,name,name_contact,mail,phone,function,code,client_id from db.hotelchain where deleted = false """
cur.execute(sql)
hotelchains= cur.fetchall()
for hotelchain in hotelchains :
    maps.hotels_chain_map.put(hotelchain[0], Hotel_Chain(hotelchain[0],hotelchain[1],hotelchain[2],hotelchain[3],hotelchain[4],hotelchain[5],hotelchain[6],hotelchain[7],hotelchain[8]))


#get hotelimages from db
#                0      1      2         3
sql = """SELECT id,hotel_id,image_name,category from db.image where deleted=false"""
cur.execute(sql)
images= cur.fetchall()
for image in images :
    maps.image_map.put(image[0], Image(image[0],image[2],image[1],image[3]))


#get hotelfacility from db
#                0              1
sql = """SELECT id_hotel,id_facility from db.hotel_facility """
cur.execute(sql)
facilities= cur.fetchall()
for facility in facilities :
    key ="0"*(8-len(str(facility[0])))+str(facility[0])+"0"*(8-len(str(facility[1])))+str(facility[1])
    sql = """SELECT name from db.facility where id= %s"""
    cur.execute(sql,(facility[1],))
    facility_name=cur.fetchone()[0]
    maps.facilities_map.put(key, Hotel_Facility(facility[0],facility[1],facility_name))


#get contracts from db
#               0   1       2       3       4       5      6        7       8       9
sql = """SELECT id,name,hotel_id,start_at,end_at,access,active,currency,market,client_id from db.contract where deleted=false  """
cur.execute(sql)
contrats= cur.fetchall()
for contrat in contrats :
    start_at=datetime.strptime(str(contrat[3]), "%Y-%m-%d").timestamp()
    end_at=datetime.strptime(str(contrat[4]), "%Y-%m-%d").timestamp()
    maps.contrats_map.put(contrat[0], Contrat(contrat[0],contrat[1],contrat[2],start_at,end_at,contrat[5],contrat[6],contrat[7],contrat[8],contrat[9]))


#get periods from db
#               0   1    2       3          4       5         6         7     8
sql = """SELECT id,name,code,contract_id,start_at,end_at,minimum_stay,delai,min_stay from db.period where deleted = false"""
cur.execute(sql)
periods= cur.fetchall()
for period in periods :
    start_at=datetime.strptime(str(period[4]), "%Y-%m-%d").timestamp()
    end_at=datetime.strptime(str(period[5]), "%Y-%m-%d").timestamp()
    maps.periods_map.put(period[0], Period(period[0],period[1],period[2],period[3],start_at,end_at,period[6],period[7],period[8]))


#get cancellation policy from db

sql = """SELECT id,periodid,boardc_id,
                 min_days_before_arrival,max_days_before_arrival,nights_to_bill
                 ,no_show,no_show_nights_to_bill, no_show_operation,
                 free_cancel, free_cancel_before,
                 room_type_id_cancel,room_id,operation,refundable FROM db.cancellationpolicy where deleted = false"""
cur.execute(sql)
cancellations= cur.fetchall()

for cancellation in cancellations :
    room_type_id_cancel = cancellation[11] if cancellation[11] is not None else 0
    room_id = cancellation[12] if cancellation[12] is not None else 0
    refundable = cancellation[14] if cancellation[14] is not None else False
    maps.cancellations_map.put(cancellation[0],
                          Cancellation(cancellation[0],cancellation[1],cancellation[2],cancellation[3],cancellation[4],cancellation[5],cancellation[6],cancellation[7],\
                          cancellation[8],cancellation[9],cancellation[10],room_type_id_cancel,room_id,cancellation[13],refundable))



#get room from db
sql = """SELECT id,name,name_code,client_id, max_pax, min_pax,child, min_adult,max_adult from db.room where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_map.put(room[0], Room(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7],room[8]))


#get room_type from db
sql = """SELECT id,name,code,client_id from db.room_type where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_type_map.put(room[0], Room_type(room[0],room[1],room[2],room[3]))


#get room_contrat from db
#               0       1        2        3             4       5           6     7     8       9           10
sql = """SELECT id,contrat_id,room_id,room_type_id,max_adult,min_adult,min_pax,max_pax,child,code_room,code_type from db.room_contrat where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_contrat_map.put(room[0], Room_contrat(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7],room[8],room[9],room[10]))



#get room_hotel from db
#               0       1        2        3
sql = """SELECT id, hotel_id, room_id, client_id FROM db.hotel_room where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_hotel_map.put(room[0], Hotel_Room(room[0],room[1],room[2],room[3]))


#get room_price from db
sql = """SELECT id,room_id,period_id,price,par_pax,extra_bed_price from db.room_prices where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_price_map.put(room[0], Room_price(room[0],room[1],room[2],room[3],room[4],room[5]))


#get room_type_price from db
sql = """SELECT id,type_id,period_id,price,par_pax,extra_bed_price from db.room_type_price where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_type_price_map.put(room[0], Room_type_price(room[0],room[1],room[2],room[3],room[4],room[5]))


#get arrangement from db
sql = """SELECT id,name,definition,client_id from db.arrangement where deleted =false"""
cur.execute(sql)
arrangements= cur.fetchall()
for arrangement in arrangements :
    maps.boarding_map.put(arrangement[0], Boarding(arrangement[0],arrangement[1],arrangement[2],arrangement[3]))

# #get arrangement_contrat from db
# #               0       1          2        3       4         5
sql = """SELECT id,contrat_id,board_id,start_at,end_at,extra_board_id from db.arrangement_contrat where deleted =false"""
cur.execute(sql)
arrangements= cur.fetchall()
for arrangement in arrangements :
    start_at=datetime.strptime(str(arrangement[3]), "%Y-%m-%d").timestamp()
    end_at=datetime.strptime(str(arrangement[4]), "%Y-%m-%d").timestamp()
    maps.board_contrat_map.put(arrangement[0], Board_contrat(arrangement[0],arrangement[1],arrangement[2],start_at,end_at,arrangement[5]))


#get board_price from db
sql = """SELECT id,arrangement_id,period_id,price,par_pax,extra_bed_price from db.arrangement_prices where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.board_price_map.put(room[0], Board_price(room[0],room[1],room[2],room[3],room[4],room[5]))


#get base_price from db
sql = """SELECT id,period_id,room_id,type_id,board_id,price,par_pax,extb_price,operation from db.base_price where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.base_price_map.put(room[0], Base_price(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7],room[8]))


#get allotement from db
sql = """SELECT id,period_id,room_id,number from db.allotement where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.allotement_map.put(room[0], Allotement(room[0],room[1],room[2],room[3]))


#get allotement_daily from db
sql = """SELECT id,contrat_id,allotement from db.allotement_daily where deleted =false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.allotement_daily_map.put(room[1], HazelcastJsonValue(room[2]))

#SpecialOffers
#get child_price from db
#               0       1       2       3       4       5       6        7
sql = """SELECT id,period_id,room_id,type_id,board_id,adults,contrat_id,childs from db.child_price where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.child_price_map.put(room[0], Child_Price(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7]))


#get special_offers from db
#               0       1               2               3           4               5          6       7        8          9        10      11          12                     13       14          15       16       17
sql = """SELECT id,receive_date,requestdate_from,requestdate_to,checkin_date,checkout_date,min_stay,max_stay,offer_type,origine,definition,note,bookingdate_range_exclusif,in_contrat,priority,contrat_id,client_id,order_number from db.special_offers where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    receive_date=datetime.strptime(str(room[1]), "%Y-%m-%d").timestamp()
    requestdate_from=datetime.strptime(str(room[2]), "%Y-%m-%d").timestamp()
    requestdate_to=datetime.strptime(str(room[3]), "%Y-%m-%d").timestamp()
    checkin_date=datetime.strptime(str(room[4]), "%Y-%m-%d").timestamp()
    checkout_date=datetime.strptime(str(room[5]), "%Y-%m-%d").timestamp()
    maps.special_offer_map.put(room[0], Special_Offer(room[0],receive_date,requestdate_from,requestdate_to,checkin_date,checkout_date,room[6],room[7],room[8],room[9],room[10],room[11],room[12],room[13],room[14],room[15],room[16],room[17]))



#get early_booking from db
#               0       1               2       3       4       5       6               7           8
# sql = """SELECT id,special_offer_id,room_id,type_id,board_id,value,value_operation,extb_value,extb_operation from db.early_booking where deleted = false"""
# cur.execute(sql)
# rooms= cur.fetchall()
# for room in rooms :
#     early_booking_map.put(room[0], Early_Booking(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7],room[8]))
# early_booking_set = early_booking_map.entry_set()
# logging.info("early_booking_set = \n {}".format(early_booking_set))

#get day_promotion from db
#               0       1               2       3       4       5        6
sql = """SELECT id,special_offer_id,rooms,room_types,boards,stay_day,pay_day from db.day_promotion where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.day_promotion_map.put(room[0], Day_Promotion(room[0],room[1],room[2],room[3],room[4],room[5],room[6]))

#get long_stay from db
#               0       1               2       3       4           5            6
# sql = """SELECT id,special_offer_id,min_days,max_days,starting_day,value,value_operation from db.long_stay where deleted = false"""
# cur.execute(sql)
# rooms= cur.fetchall()
# for room in rooms :
#     long_stay_map.put(room[0], Long_Stay(room[0],room[1],room[2],room[3],room[4],room[5],room[6]))
# long_stay_set = long_stay_map.entry_set()
# logging.info("long_stay_set = \n {}".format(long_stay_set))

#get honeymoon from db
#               0       1               2       3       4       5         6
# sql = """SELECT id,special_offer_id,room_id,type_id,board_id,value,value_operation from db.honeymoon where deleted = false"""
# cur.execute(sql)
# rooms= cur.fetchall()
# for room in rooms :
#     honeymoon_map.put(room[0], HoneyMoon(room[0],room[1],room[2],room[3],room[4],room[5],room[6]))
# honeymoon_set = honeymoon_map.entry_set()
# logging.info("honeymoon_set = \n {}".format(honeymoon_set))

#get room_promotion from db
#               0       1               2       3       4       5         6
sql = """SELECT id,special_offer_id,rooms,room_types,boards,value,value_operation from db.room_promotion where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.room_promotion_map.put(room[0], Room_Promotion(room[0],room[1],room[2],room[3],room[4],room[5],room[6]))


#get pax_promotion from db
#               0       1               2       3       4       5         6
sql = """SELECT id,special_offer_id,rooms,room_types,boards,pax_stay,pax_pay from db.pax_promotion where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.pax_promotion_map.put(room[0], Pax_Promotion(room[0],room[1],room[2],room[3],room[4],room[5],room[6]))


#get minimum_stay from db
#               0       1               2       3       4     5     6       7
# sql = """SELECT id,special_offer_id,room_id,type_id,board_id,day,start_at,end_at from db.minimum_stay where deleted = false"""
# cur.execute(sql)
# rooms= cur.fetchall()
# for room in rooms :
#     start_at=datetime.strptime(str(room[6]), "%Y-%m-%d").timestamp()
#     end_at=datetime.strptime(str(room[7]), "%Y-%m-%d").timestamp()
#     minimum_stay_map.put(room[0], Minimum_Stay(room[0],room[1],room[2],room[3],room[4],room[5],start_at,end_at))
# minimum_stay_set = minimum_stay_map.entry_set()
# logging.info("minimum_stay_set = \n {}".format(minimum_stay_set))

#get stop_sale from db
#               0       1        2       3       4     5       6        7        8
sql = """SELECT id,contrat_id,room_id,type_id,active,start_at,end_at,group_id,board_id from db.stop_sale where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    start_at=datetime.strptime(str(room[5]), "%Y-%m-%d").timestamp()
    end_at=datetime.strptime(str(room[6]), "%Y-%m-%d").timestamp()
    maps.stop_sale_map.put(room[0], Stop_Sale(room[0],room[1],room[2],room[3],room[4],room[7],start_at,end_at,room[8]))

#get weekend_day from db
#               0       1               2       3       4         5       6            7
sql = """SELECT id,contract_id,period_id,room_id,type_id,board_id,price_weekend,days from db.weekend_day where deleted = false"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    maps.weekend_day_map.put(room[0], Weekend_Day(room[0],room[1],room[2],room[3],room[4],room[5],room[6],room[7]))


#get exchange_rate from db
#               0       1             2           3      4
sql = """SELECT id,base_currency,target_currency,rate,dateandtime from db.exchange_rate"""
cur.execute(sql)
rooms= cur.fetchall()
for room in rooms :
    dateandtime =datetime.strptime(str(room[4]), "%Y-%m-%d %H:%M:%S").timestamp()
    maps.exchange_rate_map.put(room[0], Exchange_Rate(room[0],room[1],room[2],room[3],dateandtime))
logging.info("data_set End **********")


print("client shuted down")
maps.client.shutdown()
