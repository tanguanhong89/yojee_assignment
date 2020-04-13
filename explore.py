# This script is only for exploratory purposes. Hence, it is not the most optimized, nor has the best practices of
# some "data science" framework (eg. pandas), nor it has any good software engineering practices(variable names, OOPs).
# Nevertheless, it should serve its purpose well. Please read comments/documentation for better understanding
import pandas as pd
import numpy as np
from sklearn.decomposition import PCA
import plotly.graph_objects as go
import os

qwe = pd.read_csv('locations.csv', header=None)


# Non numeric  Exchange Quare
# Non numeric  Exchange Quare

def h1(v):
    try:
        return float(v)
    except:
        return False


asd = qwe.applymap(h1)  # try casting to float, else False

for c in asd:  # filters for each col with non numerical values
    asd = asd[asd[c] != False]

asd.reset_index(drop=True, inplace=True)

pcam = PCA(n_components=1)
asd1 = pcam.fit_transform(asd).reshape(-1)

ht = np.histogram(asd1, bins=3)
# bin_count = [1598    2    3] <-- 5 of them are very different from the 1598 others, assuming they are anomalies,
# remove them
# bin_interval [ -0.63874972  74.08144007 148.80162985 223.52181964]

asd2 = pd.Series(asd1)
asd = asd[asd2 < 74.08144007]  # filters out possibly anomalous data

asd.to_csv('cleaned.csv', header=None, index=False)

# Visualization. As expected, distribution is uneven.
fig = go.Figure(data=
[
    go.Scatter(x=asd[0], y=asd[1], mode='markers', marker=dict(size=4, color='green', opacity=0.8),
               name="original data"),
    go.Scatter(x=[11.552931], y=[104.933636], mode='markers', marker=dict(size=8, color='red', opacity=0.8),
               name="starting point"),
]
)
if os.path.exists(os.getcwd() + '/2.csv'):
    two = pd.read_csv('2.csv', header=None)
    # Visualization. As expected, distribution is uneven.
    fig.add_trace(go.Scatter(x=two[0], y=two[1], mode='markers', marker=dict(size=4, color='blue', opacity=0.8),
                             name="downsampled"), )

fig.update_layout(margin=dict(l=0, r=0, b=0, t=0))
fig.show()
