import { createSlice, configureStore } from '@reduxjs/toolkit'
import { persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage/session'
import { combineReducers } from 'redux'
import thunk from 'redux-thunk'
import * as category from './prpc/category_pb'

const userSlice = createSlice({
  name: 'user',
  initialState: {
    userInfo: {}
  },
  reducers: {
    setSigIn: (state, action) => {
      state.signIn = action.payload
    },
    setUserInfo: (state, action) => {
      var userInfo = action.payload
      state.userInfo = userInfo
    },
  },
})

const btSlice = createSlice({
  name: 'bt',
  initialState: {
    torrents: {},
    viodeFiles: {},
  },
  reducers: {
    updateTorrent: (state, action) => {
      var torrent = action.payload
      state.torrents[torrent.infoHash.hash] = torrent
    },
    removeTorrent: (state, action) => {
      const hash = action.payload.hash
      delete state.torrents[hash]
    },
    removeAllTorrent: (state) => {
      state.torrents = {}
    },
    updateVideoFiles: (state, action) => {
      var payload = action.payload
      state.viodeFiles[payload.infoHash.hash] = payload.btVideoMetadat
    },
  },
})

const categorySlice = createSlice({
  name: 'category',
  initialState: {
    items: {},
  },
  reducers: {
    updateItem: (state, action) => {
      let item = action.payload
      state.items[item.id] = item
    },
    deleteItem: (state, action) => {
      let itemId = action.payload
      delete state.items[itemId]
    },
  }
})

const eventSlice = createSlice({
  name: 'event',
  initialState: {
    downloadPageMouseDown: false,
  },
  reducers: {
    setDownloadPageMouse: (state, action) => {
      const isDown = action.payload
      state.downloadPageMouseDown = isDown
    },
  }
})

const persistConfig = {
  key: 'root',
  storage: storage,
}

const reducers = combineReducers({
  user: userSlice.reducer,
  bt: btSlice.reducer,
  category: categorySlice.reducer,
  event: eventSlice.reducer,
})
const persistedReducer = persistReducer(persistConfig, reducers);
const store = configureStore({
  reducer: persistedReducer,
  devTools: process.env.NODE_ENV !== 'production',
  middleware: [thunk]
})

const selectUserInfo = (state) => {
  return state.user.userInfo
}

const selectTorrent = (state, infoHash) => {
  return state.bt.torrents[infoHash.hash]
}

const selectInfoHashs = (state) => {
  const keys = []
  Object.values(state.bt.torrents).forEach((v) => {
    var infoHash = {
      version: v.infoHash.version,
      hash: v.infoHash.hash
    }
    keys.push(infoHash)
  })
  return keys
}

const selectBtVideoFiles = (state, infoHash) => {
  return state.bt.viodeFiles[infoHash.hash]
}

const selectCategoryItem = (state, itemId) => {
  return state.category.items[itemId]
}

const selectCategorySubItems = (state, parentId) => {
  const ds = []
  if (!state.category.items[parentId]) {
    return ds
  }
  let sudItemIds = state.category.items[parentId].subItemIdsList
  if (!sudItemIds) {
    return ds
  }
  sudItemIds.map((id) => {
    if (state.category.items[id]) {
      ds.push(state.category.items[id])
    }
    return null
  })
  return ds
}

const selectSubDirectory = (state, parentId) => {
  const ds = []
  if (!state.category.items[parentId]) {
    return ds
  }
  let sudItemIds = state.category.items[parentId].subItemIdsList
  if (!sudItemIds) {
    return ds
  }
  sudItemIds.map((id) => {
    if (state.category.items[id] &&
      state.category.items[id].typeId === category.CategoryItem.Type.DIRECTORY) {
      ds.push(state.category.items[id])
    }
    return null
  })
  return ds
}

const isDownloadPageMouseDown = (state) => {
  return state.event.downloadPageMouseDown
}

export {
  store, userSlice, btSlice, categorySlice, eventSlice,
  selectUserInfo,
  selectTorrent, selectInfoHashs, selectBtVideoFiles,
  selectCategoryItem, selectCategorySubItems, selectSubDirectory,
  isDownloadPageMouseDown
}