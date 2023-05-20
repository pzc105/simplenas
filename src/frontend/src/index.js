import React from 'react';
import { Provider } from 'react-redux';
import ReactDOM from 'react-dom/client';
import APP from './router.js'
import { store } from './store.js'
import { PersistGate } from 'redux-persist/integration/react'
import { persistStore } from 'redux-persist'

let persistor = persistStore(store);
const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <Provider store={store}>
    <PersistGate loading={null} persistor={persistor}>
      <APP />
    </PersistGate>
  </Provider>
)