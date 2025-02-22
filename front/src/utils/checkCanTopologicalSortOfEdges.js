/**
 * 检查联表条件数组可否进行拓扑排序
 * @param {*} edges  联表条件数组
 * @returns  true: 能够拓扑排序， false: 不能拓扑排序
 */
const checkCanTopologicalSortOfEdges = (edges) => {
    const adjList = {};
    const inDegree = {};
    const allpos = {};

    // Initialize graph
    for (const edge of edges) {
        for (const pos of [edge.pos1, edge.pos2]) {
            if (!allpos[pos]) {
                adjList[pos] = [];
                allpos[pos] = true;
                inDegree[pos] = 0;
            }
        }
    }

    // Build graph
    for (const edge of edges) {
        adjList[edge.pos1].push(edge.pos2);
        adjList[edge.pos2].push(edge.pos1);
        inDegree[edge.pos1]++;
        inDegree[edge.pos2]++;
    }

    // Initialize queue with nodes having 0 in-degree
    const queue = [];
    for (const pos in inDegree) {
        if (inDegree[pos] === 1) { // Adjusted for undirected graph
            queue.push(pos);
        }
    }

    // Process nodes in the queue
    const sortedEdges = [];
    const visited = {};
    while (queue.length > 0) {
        const pos = queue.shift();

        // Find the edge that should be added to the sorted list
        for (const edge of edges) {
            if ((edge.pos1 === pos || edge.pos2 === pos) &&
                !visited[edge.pos1 + edge.pos2] &&
                !visited[edge.pos2 + edge.pos1]) {
                sortedEdges.push(edge);
                visited[edge.pos1 + edge.pos2] = true;
                visited[edge.pos2 + edge.pos1] = true;
                break;
            }
        }

        // Decrease in-degree of adjacent nodes
        for (const adj of adjList[pos]) {
            inDegree[adj]--;
            if (inDegree[adj] === 1) { // Adjusted for undirected graph
                queue.push(adj);
            }
        }
    }

    // Check for cycles in the graph
    if (sortedEdges.length !== edges.length) {
        return false;
    }


    console.log(sortedEdges);

    return true;
}

export default checkCanTopologicalSortOfEdges;