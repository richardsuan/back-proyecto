from dash import dcc

def render():
    return dcc.Tabs(
        id='graph-tabs',
        value='Volumen',
        className='custom-tabs',
        children=[
            dcc.Tab(label='Volumen', value='Volumen', className='custom-tab', selected_className='custom-tab--selected'),
            dcc.Tab(label='Presion', value='Presion', className='custom-tab', selected_className='custom-tab--selected'),
            dcc.Tab(label='Temperatura', value='Temperatura', className='custom-tab', selected_className='custom-tab--selected'),
        ]
    )