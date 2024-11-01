import useSWR from "swr";
import PostForm from "../components/PostForm";
import Datafeed from "../components/Datafeed";

const fetcher = url => fetch(url).then(res => res.json());

export default function Home() {
    const { data: posts, mutate } = useSWR("http://localhost:8080/api/posts", fetcher, { refreshInterval: 5000 });

    return (
        <div>
            <h1>Datafeed</h1>
            <PostForm onPostSubmit={mutate} />
            <Datafeed posts={posts || []} />
        </div>
    );
}
