import sys
import json
import time
import random
import urllib.request
from threading import Thread, Lock

'''
可调整的参数
BEGIN
'''

domain = 'http://localhost:8080'

urls = [
        domain+'/ogm/approval/Healthy/Echo',
        ]

params = {
        'msg':'hello',
        }

options = {
        'content-type':'application/json;charset=UTF-8'
        }

thread_count = 20
invoke_count = 20
interval = 0


'''
可调整的参数
END
'''



failed = 0
successed = 0
thread_finished = 0
errors = {}
cps = 0 #call per seconds
cps_successed = 0
cps_failed = 0
exit = False
def onThreadFinish(_id):
    global thread_finished
    lock.acquire()
    thread_finished = thread_finished + 1
    lock.release()

def call(_id, _n):
    global failed
    global successed
    global cps_failed
    global cps_successed
    global interval
    global errors
    global cps
    for i in range(0, _n):
        time.sleep(interval)
        data = json.dumps(params)
        data = bytes(data, 'utf8')
        url = urls[random.randint(0, len(urls)-1)]
        request = urllib.request.Request(url, data, options)
        try:
            req = urllib.request.urlopen(request)
            reply = req.read().decode('utf-8')
            req.close()
            del req
            #print(reply)
            lock.acquire()
            successed = successed + 1
            cps_successed = cps_successed + 1
            lock.release()
        except Exception as e:
            err = str.format('%s -> %s'%(url,e))
            lock.acquire()
            failed = failed + 1
            cps_failed = cps_failed + 1
            errors[err] = 0
            lock.release()
        lock.acquire()
        cps = cps + 1
        lock.release()
    onThreadFinish(_id)

def tick_print():
    global cps
    global cps_successed
    global cps_failed
    tick = 0
    while True:
        time.sleep(1)
        tick = tick + 1
        print('> 【每秒】访问次数：%d, 成功：%d，失败：%d, 【总计】耗时：%d秒，成功：%d，失败：%d，完成线程：%d，剩余任务：%d'%(cps, cps_successed, cps_failed, tick, successed, failed, thread_finished, thread_count*invoke_count - failed-successed))
        lock.acquire()
        cps = 0
        cps_successed = 0
        cps_failed = 0
        lock.release()
        if thread_finished == thread_count:
            break
    print('----------------------------------------------------------')
    print('错误列表：')
    for err in errors:
        print(err)
    print('##########################################################')

print('##########################################################')
print('接口列表：')
for url in urls:
    print(url)
print('----------------------------------------------------------')
print('线程数：%d '%(thread_count))
print('每线程访问次数：%d '%(invoke_count))
print('访问间隔时间：%f '%(interval))
print('**********************************************************')
pt = Thread(target=tick_print)
pt.start()
# 创建100个线程
for i in range(0,thread_count):
    lock = Lock()
    # 每个线程10次访问
    t = Thread(target=call, args=(i, invoke_count,))
    t.start()
