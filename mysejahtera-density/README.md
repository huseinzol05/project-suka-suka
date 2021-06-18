# MySejahtera Density

Generate resolution points using meshgrid and request each points to MySejahtera Hotspot API to get active cases in 1KM radius.

## how-to

1. Run [generate-malaysia-points.ipynb](generate-malaysia-points.ipynb) to generate resolution points for each states. You might want to finetune `resolution` variable to define better definition of resolution points,

```python
resolutions = {'Sarawak': 0.005, 'Sabah': 0.005}
default_resolution = 0.004
```

These resolution points are pretty naive, not considering rivers and non-walkable terrains.

2. Run [get-malaysia-points-case.ipynb](get-malaysia-points-case.ipynb) to request each resolution points to MySejahtera Hotspot API.

```python
resolution_points = glob('data/*-points.json')
resolution_points_with_case = glob('data/*-points.json-points-cases.json')
```

I will run this notebook occasionally, just check last status update at [last-update.json](last-update.json).

## Why

1. We want the public able to do analysis for the active cases.
