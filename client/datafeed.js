export default function Datafeed({ posts }) {
    return (
        <div>
            {posts.map(post => (
                <div key={post.id} className="post">
                    <h3>{post.username}</h3>
                    <p>{post.content}</p>
                    <img src={`http://localhost:8080/${post.image_url}`} alt="Post image" width={300} />
                    <small>{new Date(post.timestamp).toLocaleString()}</small>
                </div>
            ))}
        </div>
    );
}
