import { Routes, Route, Navigate, BrowserRouter } from 'react-router-dom';
import LoginPage from './pages/LoginPage';
import RegisterPage from './pages/RegisterPage';
import { useAppSelector } from './hooks/redux';
import TopicsPage from './pages/TopicsPage';
import TopicPostsPage from './pages/TopicPostsPage';
import PostPage from './pages/PostPage';
import CreatePostPage from './pages/CreatePostPage';
import Layout from './components/Layout';
import EditPostPage from './pages/EditPostPage';
import CreateTopicPage from './pages/CreateTopicPage';
import ProfilePage from './pages/ProfilePage';
import EditTopicPage from './pages/EditTopicPage';
import { ThemeProvider } from './contexts/ThemeContext';

function App() {
  const { isAuthenticated } = useAppSelector(state => state.auth);

  return (
    <ThemeProvider>
      <BrowserRouter>
        <Routes>
          {/* Public Routes (No Layout) */}
          <Route path='/login' element={<LoginPage />} />
          <Route path='/register' element={<RegisterPage />} />

          {/* Protected Routes (With Layout) */}
          <Route 
            path='/profile'
            element={
              isAuthenticated 
                ? (<Layout>
                    <ProfilePage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route 
            path='/users/:userID'
            element={
              isAuthenticated 
                ? (<Layout>
                    <ProfilePage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route 
            path='/topics' 
            element={
              isAuthenticated 
                ? (<Layout>
                    <TopicsPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            } 
          />

          <Route 
            path='/topics/create' 
            element={
              isAuthenticated 
                ? (<Layout>
                    <CreateTopicPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route 
            path='/topics/:topicID' 
            element={
              isAuthenticated 
                ? (<Layout>
                    <TopicPostsPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            } 
          />

          <Route
            path='/topics/:topicID/edit'
            element={
              isAuthenticated 
                ? (<Layout>
                    <EditTopicPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route 
            path='/topics/:topicID/create-post'
            element={
              isAuthenticated 
                ? (<Layout>
                    <CreatePostPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route
            path='/posts/:postID'
            element={
              isAuthenticated 
                ? (<Layout>
                    <PostPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          <Route
            path='/posts/:postID/edit'
            element={
              isAuthenticated 
                ? (<Layout>
                    <EditPostPage />
                  </Layout>) 
                : (<Navigate to="/login" replace />)
            }
          />

          {/* Default Route */}
          <Route 
            path='/' 
            element={isAuthenticated ? <Navigate to="/topics" replace /> : <Navigate to="/login" replace />} 
          />
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;