# http://stackoverflow.com/questions/17756925/how-to-plot-heatmap-colors-in-3d-in-matplotlib

import json

import numpy as np
from pylab import *
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt

with open('net.json') as data_file:    
    nodes = json.load(data_file)['nodes']

def randrange(n, vmin, vmax):
    return (vmax-vmin)*np.random.rand(n) + vmin

fig = plt.figure(figsize=(8,6))

ax = fig.add_subplot(111,projection='3d')
n = 100

# xs = randrange(n, 0, 100)
# ys = randrange(n, 0, 100)
# zs = randrange(n, 0, 100)
# the_fourth_dimension = randrange(n,0,100)
xs = np.array([])
ys = np.array([])
zs = np.array([])
values = []
for node in nodes:
    xs = np.append(xs, node['position'][0])
    ys = np.append(ys, node['position'][1])
    zs = np.append(zs, node['position'][2])
    values = np.append(values, node['value'])

colors = cm.hsv(values/max(values))

colmap = cm.ScalarMappable(cmap=cm.hsv)
colmap.set_array(values)

yg = ax.scatter(xs, ys, zs, c=colors, alpha=0.5, marker='.') # or set marker to "o"
cb = fig.colorbar(colmap)

ax.set_xlabel('X')
ax.set_ylabel('Y')
ax.set_zlabel('Z')


plt.savefig('net.png')