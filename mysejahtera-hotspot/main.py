import json
import pandas as pd
import geopandas
from shapely.geometry.multilinestring import MultiLineString
from shapely.geometry.collection import GeometryCollection
from shapely.geometry.multipolygon import MultiPolygon
from shapely.ops import cascaded_union
from shapely.geometry.polygon import Polygon
from shapely.geometry import Point
from geopandas import GeoDataFrame
from libpysal.cg.alpha_shapes import alpha_shape_auto
from sklearn.cluster import DBSCAN
from datetime import datetime, date
import urllib.parse
import libpysal as lps
import numpy as np
import os
import copy

today = str(date.today())
today

STATES = {
    'johor': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Johor-points.json-points-cases.json',
    'kedah': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Kedah-points.json-points-cases.json',
    'kelantan': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Kelantan-points.json-points-cases.json',
    'kl': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Federal%20Territory%20of%20Kuala%20Lumpur-points.json-points-cases.json',
    'labuan': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Labuan-points.json-points-cases.json',
    'melaka': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Melaka-points.json-points-cases.json',
    'negeri-sembilan': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Negeri%20Sembilan-points.json-points-cases.json',
    'pahang': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Pahang-points.json-points-cases.json',
    'penang': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Penang-points.json-points-cases.json',
    'perak': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Perak-points.json-points-cases.json',
    'putrajaya': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Federal%20Territory%20of%20Putrajaya-points.json-points-cases.json',
    'sabah': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Sabah-points.json-points-cases.json',
    'sarawak': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Sarawak-points.json-points-cases.json',
    'selangor': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Selangor-points.json-points-cases.json',
    'terengganu': 'https://raw.githubusercontent.com/huseinzol05/project-suka-suka/main/mysejahtera-density/data/Terengganu-points.json-points-cases.json',
}

COLOR = {0: '#0000FF', 1: '#0000FF', 2: '#0000FF', 3: '#0000FF',
         4: '#0000FF', 5: '#0000FF', 6: '#0000FF', 7: '#FFFF00',
         8: '#FFFF00', 9: '#FFFF00', 10: '#FFFF00', 11: '#FFFF00',
         12: '#FFFF00', 13: '#FF0000', 14: '#FF0000', 15: '#FF0000',
         16: '#FF0000', 17: '#FF0000', 18: '#FF0000', 19: '#FF0000'}
COLORS = ['#0000FF', '#FFFF00', '#FF0000']
STEP = 0.05
DEFAULT_DISTANCE = 0.0045
DISTANCES = {'sabah': 0.0075, 'sarawak': 0.0085}


def get_cluster_boundary(labels, xys, scores, xy=["X", "Y"], crs=None, step=1):
    try:
        from geopandas import GeoSeries
    except ModuleNotFoundError:

        def GeoSeries(data, index=None, crs=None):
            return list(data)

    lbl_type = type(labels.iloc[0])
    noise = lbl_type(-1)
    ids_in_cluster = labels[labels != noise].index
    g = xys.loc[ids_in_cluster, xy].groupby(labels[ids_in_cluster])
    polys = []
    cluster_lbls = []
    y, totals = [], []
    highests = []
    for sub in g.groups:
        try:
            max_val = np.percentile(xys.loc[g.groups[sub], 2], 80)
            highest = xys.loc[g.groups[sub], [0, 1, 2]]
            highest = highest[highest[2] >= max_val]
            highest = np.mean(highest, 0).tolist()
            polys.append(_asa((xys.loc[g.groups[sub], xy].values, 1)))
            y.append(scores.loc[g.groups[sub]].mean())
            totals.append(scores.loc[g.groups[sub]].shape[0])
            cluster_lbls.append(sub)
            highests.append(highest)
        except:
            pass
    polys = GeoSeries(polys, index=cluster_lbls, crs=crs)
    return polys, y, totals, highests


def _asa(pts_s):
    return alpha_shape_auto(pts_s[0], step=pts_s[1])


def check_boundaries(v):
    for i in range(len(boundaries)):
        if boundaries[i][0] <= v < boundaries[i][1]:
            return i


for STATE, LINK in STATES.items():
    print(STATE, LINK)
    file = urllib.parse.unquote(LINK.split('/')[-1])
    if not os.path.exists(file):
        os.system(f'wget {LINK}')

    with open(file) as fopen:
        data = json.load(fopen)
    df = pd.DataFrame(data)
    nonzero = df[2][df[2] > 0]
    arange = np.arange(0, 1.0 + STEP, step=STEP)
    percentiles = np.percentile(nonzero, arange * 100)
    boundaries = [percentiles[i: i + 2].tolist() for i in range(len(percentiles))][:-1]
    boundaries = [[int(i[0]), int(i[1])] for i in boundaries]
    boundaries[-1][-1] += 1
    geometry = [Point(xy) for xy in zip(df[0], df[1])]
    crs = {'init': 'epsg:4326'}
    bl_gdf = GeoDataFrame(df, crs=crs, geometry=geopandas.points_from_xy(df[1], df[0]))
    db = bl_gdf.to_crs(epsg=5243)
    db["X"] = db.geometry.x
    db["Y"] = db.geometry.y

    already_processed = set()
    polygons = []
    for i in range(len(boundaries)):
        print(i, boundaries[i])

        if (boundaries[i][1] - boundaries[i][0]) < 1:
            continue

        if str(boundaries[i]) in already_processed:
            continue

        try:
            dbscan = DBSCAN(eps=DISTANCES.get(STATE, DEFAULT_DISTANCE), min_samples=3)
            filtered_df = df[(df[2] >= boundaries[i][0]) & (df[2] < boundaries[i][1])]
            filtered_df_index = filtered_df.index
            clustering = dbscan.fit(filtered_df[[0, 1]].values)

            labels = np.full(shape=df.shape[0], fill_value=-1)
            for no in range(len(clustering.labels_)):
                labels[filtered_df_index[no]] = clustering.labels_[no]

            polys, ys, totals, highests = get_cluster_boundary(pd.Series(labels), db, db[2], crs=db.crs)
            polys = polys.to_crs('crs')

            for k in range(len(polys)):
                if polys.iloc[k].area <= 1e-12:
                    continue
                polygons.append({'polygon': polys.iloc[k].convex_hull, 'y': [ys[k]],
                                'total': totals[k], 'color': COLOR[i], 'highest': [highests[k]]})

            already_processed.add(str(boundaries))
        except Exception as e:
            print(e)

    polygons_ = copy.deepcopy(polygons)

    r = []
    processed = set()
    for i in reversed(range(len(polygons_))):
        if i in processed:
            continue

        for k in reversed(range(len(polygons_))):
            if k in processed:
                continue

            if i == k:
                continue
            if k > i:
                continue

            if not polygons_[i]['polygon'].intersects(polygons_[k]['polygon']):
                continue

            l = polygons_[i]['polygon'].union(polygons_[k]['polygon'])
            if not isinstance(l, MultiPolygon):
                if isinstance(l, GeometryCollection):
                    l = [p for p in l if isinstance(p, Polygon)]
                    l = cascaded_union(l)
                polygons_[i]['polygon'] = l.convex_hull
                polygons_[i]['y'].extend(polygons_[k]['y'])
                polygons_[i]['highest'].extend(polygons_[k]['highest'])
                polygons_[i]['total'] += polygons_[k]['total']
                processed.add(k)

        polygons_[i]['y'] = np.mean(polygons_[i]['y'])
        new_color = COLOR[check_boundaries(polygons_[i]['y'])]
        polygons_[i]['color'] = new_color
        if polygons_[i]['polygon'].area > 1e-12:
            r.append(polygons_[i])

    post = []
    for i in reversed(range(len(r))):

        inside = False
        a = r[i]['polygon'].area

        for k in range(len(r)):
            if k in processed:
                continue

            if i == k:
                continue

            if i < k:
                continue

            if not r[i]['polygon'].intersects(r[k]['polygon']):
                continue

            l = r[i]['polygon'].intersection(r[k]['polygon'])
            if not isinstance(l, MultiPolygon):
                if isinstance(l, GeometryCollection):
                    l = [p for p in l if isinstance(p, Polygon)]
                    l = cascaded_union(l)
                if l.area / a >= 0.8:
                    inside = True
                    break

        if not inside:
            post.append(r[i])

    for i in range(len(post)):
        polygons_ = []
        x, y = post[i]['polygon'].exterior.coords.xy
        area = post[i]['polygon'].area
        for x_, y_ in zip(x, y):
            polygons_.append({'lat': y_, 'lng': x_})
        post[i]['polygon'] = polygons_
        post[i]['area'] = area

    with open(f'data/{STATE}.json', 'w') as fopen:
        json.dump(post, fopen)

    os.remove(file)

with open('last-update.json', 'w') as fopen:
    json.dump({'last-update': today}, fopen)
