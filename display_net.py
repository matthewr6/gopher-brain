# http://stackoverflow.com/questions/17756925/how-to-plot-heatmap-colors-in-3d-in-matplotlib

import sys
import json

import numpy as np
from pylab import *
from mpl_toolkits.mplot3d import Axes3D
import matplotlib.pyplot as plt

COLOR_MAP = cm.plasma

for num in range(0, int(sys.argv[1])):
    with open('frames/net_{num}.json'.format(num=num)) as data_file:   
        nodes = json.load(data_file)['nodes']

    def randrange(n, vmin, vmax):
        return (vmax-vmin)*np.random.rand(n) + vmin

    fig = plt.figure(figsize=(8,6))

    ax = fig.add_subplot(111,projection='3d')
    n = 100

    xs = np.array([])
    ys = np.array([])
    zs = np.array([])
    values = []
    for node in nodes:
        xs = np.append(xs, node['position'][0])
        ys = np.append(ys, node['position'][1])
        zs = np.append(zs, node['position'][2])
        values = np.append(values, node['value'])

    colors = COLOR_MAP(values/max(values))

    colmap = cm.ScalarMappable(cmap=COLOR_MAP)
    colmap.set_array(values)

    yg = ax.scatter(xs, ys, zs, c=colors, alpha=0.5, marker='.') # or set marker to "o"
    cb = fig.colorbar(colmap)

    ax.set_xlabel('X')
    ax.set_ylabel('Y')
    ax.set_zlabel('Z')


    plt.savefig('frames/net_{num}.png'.format(num=num))
    # plt.close()