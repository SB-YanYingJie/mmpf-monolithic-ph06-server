# mpstat.py
import csv

cols = ('time',
  'centos7-1','centos7-2')

with open('../sh/output.csv') as f:
  with open('./container.mesurement.memory.result.csv', mode='w') as g:
    reader = csv.reader(f)
    writer = csv.DictWriter(g, fieldnames=cols, lineterminator='\n')

    r = {'time': ''}
    
    for row in reader:
      time, c_name, _, mem ,*_ = row
      print(time,c_name,mem)
      if r['time'] != time:
        if r['time'] != '':
          writer.writerow(r)
        r = {'time': time}

      r[c_name] = mem
    writer.writerow(r)
    
