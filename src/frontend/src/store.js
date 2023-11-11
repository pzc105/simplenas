import { createSlice, configureStore } from '@reduxjs/toolkit'
import { persistReducer } from 'redux-persist';
import storage from 'redux-persist/lib/storage'
import { combineReducers } from 'redux'
import thunk from 'redux-thunk'
import * as category from './prpc/category_pb'

const userSlice = createSlice({
  name: 'user',
  initialState: {
    userInfo: null,
    shownChatPanel: false,
    openGlobalChat: false,
    globalChatPosition: { 0: { x: null, y: null } },
  },
  reducers: {
    setUserInfo: (state, action) => {
      var userInfo = action.payload
      state.userInfo = userInfo
    },
    setShowChatPanel: (state, action) => {
      state.shownChatPanel = action.payload
    },
    setOpenGlobalChat: (state, action) => {
      state.openGlobalChat = action.payload
    },
    setGlobalChatPosition: (state, action) => {
      state.globalChatPosition[action.payload["itemId"]] = action.payload
    },
  },
})

const btSlice = createSlice({
  name: 'bt',
  initialState: {
    torrentInfo: {},
    torrentStatus: {},
    viodeFiles: {},
  },
  reducers: {
    updateTorrentInfo: (state, action) => {
      var torrent = action.payload
      state.torrentInfo[torrent.infoHash.hash] = torrent
    },
    updateTorrentStatus: (state, action) => {
      var torrent = action.payload
      state.torrentStatus[torrent.infoHash.hash] = torrent
    },
    removeTorrent: (state, action) => {
      const hash = action.payload.hash
      delete state.torrentStatus[hash]
    },
    removeAllTorrent: (state) => {
      state.torrentStatus = {}
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
    videoInfos: {},
  },
  reducers: {
    clear: (state) => {
      state.items = {}
      state.videoInfos = {}
    },
    updateItem: (state, action) => {
      let item = action.payload
      state.items[item.id] = item
    },
    deleteItem: (state, action) => {
      let itemId = action.payload
      delete state.items[itemId]
    },
    updateVideoInfo: (state, action) => {
      state.videoInfos[action.payload.itemId] = action.payload.videoInfo
    },
  }
})

const magnetShares = createSlice({
  name: 'magnetShares',
  initialState: {
    parentId: 0,
    magnetSharesItems: [],
  },
  reducers: {
    updateMagnetSharesItems: (state, action) => {
      let items = action.payload.items
      let parentId = action.payload.parentId
      state.parentId = parentId
      state.magnetSharesItems = items
    },
  }
})

const playerSlice = createSlice({
  name: 'player',
  initialState: {
    selectedAudio: {},
    autoPlayVideo: false,
  },
  reducers: {
    updateSelectedAudio: (state, action) => {
      state.selectedAudio[action.payload.vid] = action.payload.aid
    },
    setAutoContinuedPlayVideo: (state, action) => {
      state.autoPlayVideo = action.payload
    }
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
  player: playerSlice.reducer,
  event: eventSlice.reducer,
  magnetShares: magnetShares.reducer,
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

const selectShownChatPanel = (state) => {
  return state.user.shownChatPanel
}

const selectOpenGlobalChat = (state) => {
  return state.user.openGlobalChat
}

const selectGlobalChatPosition = (state, itemId) => {
  if (state.user.globalChatPosition) {
    return state.user.globalChatPosition[itemId]
  }
}

const selectTorrent = (state, infoHash) => {
  if (state.bt.torrentInfo) {
    return state.bt.torrentInfo[infoHash.hash]
  }
}

const selectTorrentStatus = (state, infoHash) => {
  if (state.bt.torrentStatus) {
    return state.bt.torrentStatus[infoHash.hash]
  }
}

const selectInfoHashs = (state) => {
  const keys = []
  Object.values(state.bt.torrentStatus).forEach((v) => {
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

const selectCategoryItems = (state, ...itemIds) => {
  let ret = {}
  for (let itemId of itemIds) {
    if (state.category.items[itemId]) {
      ret[itemId] = state.category.items[itemId]
    }
  }
  return ret
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

const selectMagnetSharesItems = (state) => {
  return state.magnetShares.magnetSharesItems
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

const selectItemVideoInfo = (state, itemId) => {
  return state.category.videoInfos[itemId]
}

const isDownloadPageMouseDown = (state) => {
  return state.event.downloadPageMouseDown
}

const getSelectedAudio = (state, vid) => {
  if (vid in state.player.selectedAudio) {
    return state.player.selectedAudio[vid]
  }
  return null
}

const selectAutoPlayVideo = (state) => {
  return state.player.autoPlayVideo
}

export {
  store, userSlice, btSlice, categorySlice, eventSlice, playerSlice, magnetShares,
  selectUserInfo, selectShownChatPanel, selectOpenGlobalChat, selectGlobalChatPosition,
  selectTorrent, selectInfoHashs, selectBtVideoFiles, selectTorrentStatus,
  selectCategoryItem, selectCategoryItems, selectCategorySubItems, selectSubDirectory, selectItemVideoInfo,
  selectMagnetSharesItems,
  getSelectedAudio, selectAutoPlayVideo,
  isDownloadPageMouseDown
}