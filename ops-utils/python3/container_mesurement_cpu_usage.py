# mpstat.py
import csv

cols = ('time',
  'mmpf-monolithic_devcontainer-app-1', 'mmpf-monolithic_devcontainer-redis-1')

with open('../sh/output.csv') as f:
  with open('./container.mesurement.cpu.result.csv', mode='w') as g:
    reader = csv.reader(f)
    writer = csv.DictWriter(g, fieldnames=cols, lineterminator='\n')

    r = {'time': ''}
    
    for row in reader:
      time, c_name, cpu , *_ = row
      print(time,c_name,cpu)
      if r['time'] != time:
        if r['time'] != '':
          writer.writerow(r)
        r = {'time': time}

      r[c_name] = cpu
    writer.writerow(r)
    
