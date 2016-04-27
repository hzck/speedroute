var app = angular.module('speedrunRouting', ['ngVis', 'ui.bootstrap']);

app.controller('RouteCtrl', function($log, $http, VisDataSet, $location) {

    var g = this;

    //create page
    g.name
    g.password
    g.livesplit

    g.rewards = [];
    g.nodes = []
    g.edges = [];
    g.startNode = "";
    g.endNode = "";
    g.shortestPath = [];

    var rewardBeingEdited = undefined;
    var nodeBeingEdited = undefined;
    var edgeBeingEdited = undefined;

    /* vis network data start */
    var networkNodes = new VisDataSet();
    var networkEdges = new VisDataSet();
    g.network = {
        nodes: networkNodes,
        edges: networkEdges
    };
    g.options = {
        //autoResize: true
        height: '100%'
    };
    g.events = {};
    /* vis network data end */

    g.resetReward = function() {
        rewardBeingEdited = undefined;
        g.reward = {
            edit: false,
            id: "",
            unique: false
        }
    };

    resetRewardRef = function(obj) {
        obj.rewardRef = {
            edit: false,
            rewardId: "",
            quantity: ""
        }
    };

    g.resetNode = function() {
        nodeBeingEdited = undefined;
        g.node = {
            edit: false,
            id: "",
            revisitable: false,
            rewards: []
        }
        resetRewardRef(g.node);
    };

    resetEdgeWeight = function() {
        g.edge.weight = {
            edit: false,
            time: "",
            requirements: []
        }
        resetRewardRef(g.edge.weight);
    };

    g.resetEdge = function() {
        edgeBeingEdited = undefined;
        g.edge = {
            edit: false,
            from: "",
            to: "",
            weights: []
        }
        resetEdgeWeight();
    };

    containsEdge = function(list, from, to) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].from === from && list[i].to === to) {
                return true;
            }
        }
        return false;
    }

    contains = function(list, id) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].id === id) {
                return true;
            }
        }
        return false;
    };

    containsReward = function(list, id) {
        for(var i = 0; i < list.length; i++) {
            if(list[i].rewardId === id) {
                return true;
            }
        }
        return false;
    };

    removeObj = function(list, obj) {
        var index = list.indexOf(obj);
        if(index > -1) {
            list.splice(index, 1);
        }
    }

    canRewardBeRemoved = function(id) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            for(var j = 0; j < edge.weights.length; j++) {
                var weight = edge.weights[j];
                for(var k = 0; k < weight.requirements.length; k++) {
                    if(weight.requirements[k].rewardId === id) {
                        return false;
                    }
                }
            }
        }
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            for(var j = 0; j < node.rewards.length; j++) {
                if(node.rewards[j].rewardId === id) {
                    return false;
                }
            }
        }
        return true;
    }

    canNodeBeRemoved = function(id) {
        for(var i = 0; i < g.edges.length; i++) {
            if(g.edges[i].from === id || g.edges[i].to === id) {
                return false;
            }
        }
        return true;
    }

    toggleEdit = function(reward, node, edge) {
        g.reward.edit = reward;
        g.node.edit = node;
        g.edge.edit = edge;
    };

    log = function(msg) {
        $log.debug(msg);
    };

    updateNodeRewardReferences = function(oldId, newId) {
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            for(var j = 0; j < node.rewards.length; j++) {
                var nodeRew = node.rewards[j];
                if(nodeRew.rewardId === oldId) {
                    nodeRew.rewardId = newId;
                    log("updated node reward reference");
                }
            }
        }
    };

    updateEdgeRequirementReferences = function(oldId, newId) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            for(var j = 0; j < edge.weights.length; j++) {
                var weight = edge.weights[j];
                for(var k = 0; k < weight.requirements.length; k++) {
                    var weightReq = weight.requirements[k];
                    if(weightReq.rewardId === oldId) {
                        weightReq.rewardId = newId;
                        log("updated edge weight requirement reference");
                    }
                }
            }
        }
    };

    updateEdgeNodesReferences = function(oldId, newId) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            if(edge.from === oldId) {
                edge.from = newId;
                log("updated edge from node reference");
            }
            if(edge.to === oldId) {
                edge.to = newId;
                log("updated edge to node reference");
            }
        }
    };

    getNodeIndex = function(id) {
        for(var i = 0; i < g.nodes.length; i++) {
            var node = g.nodes[i];
            if(node.id === id) {
                return i;
            }
        }
        return -1;
    };

    getEdgeIndex = function(from, to) {
        for(var i = 0; i < g.edges.length; i++) {
            var edge = g.edges[i];
            if(edge.from === from && edge.to === to) {
                return i;
            }
        }
        return -1;
    };

    g.len = function(list) {
        if (list) {
            return list.length;
        }
        return 0;
    };

    g.addReward = function() {
        var id = g.reward.id;
        if(id) {
            if(rewardBeingEdited) {
                log("EDIT MODE");
                var oldId = rewardBeingEdited.id;
                var isNewIdOk = oldId !== id && !contains(g.rewards, id);
                if(oldId === id || isNewIdOk) {
                    rewardBeingEdited.unique = g.reward.unique;
                    if(isNewIdOk) {
                        rewardBeingEdited.id = id;
                        updateNodeRewardReferences(oldId, id);
                        updateEdgeRequirementReferences(oldId, id);
                    }
                    g.resetReward();
                } else {
                    log("failed to update reward");
                }
            } else {
                log("ADD MODE");
                if(!contains(g.rewards, id)) {
                    g.rewards.push({
                        id: id,
                        unique: g.reward.unique
                    });
                    g.resetReward();
                } else {
                    log("failed to add reward");
                }
            }
        } else {
            log("g.reward.id is not set");
        }
    };

    g.addNodeReward = function() {
        var id = g.node.rewardRef.rewardId;
        if(id && !containsReward(g.node.rewards, id) && contains(g.rewards, id)) {
            g.node.rewards.push({
                rewardId: id,
                quantity: (parseInt(g.node.rewardRef.quantity) || 1)//check if this can be removed if not set
            });
            resetRewardRef(g.node);
        }
    }

    g.addNode = function() {
        var id = g.node.id;
        if(id && !g.node.rewardRef.edit) {
            if(nodeBeingEdited) {
                log("EDIT MODE");
                var oldId = nodeBeingEdited.id;
                var isNewIdOk = oldId !== id && !contains(g.nodes, id);
                if(oldId === id || isNewIdOk) {
                    nodeBeingEdited.revisitable = g.node.revisitable;
                    nodeBeingEdited.rewards = g.node.rewards;
                    if(isNewIdOk) {
                        log("newId ok");
                        nodeBeingEdited.id = id;
                        updateEdgeNodesReferences(oldId, id);
                        networkNodes.update({
                            id: g.nodes.indexOf(nodeBeingEdited),
                            label: id,
                            color: 'lightblue'
                        });
                    }
                    g.resetNode();
                } else {
                    log("failed to update node");
                }
            } else {
                log("ADD MODE");
                if (!contains(g.nodes, id)) {
                    g.nodes.push({
                        id: id,
                        revisitable: g.node.revisitable,
                        rewards: g.node.rewards
                    });
                    networkNodes.add({
                        id: g.nodes.length-1,
                        label: id,
                        color: 'lightblue'
                    });
                    g.resetNode();
                } else {
                    log("failed to add node");
                }
            }
        } else {
            log("g.node.id is not set or node.rewards are being edited");
        }
    };

    g.addEdgeRequirement = function() {
        var id = g.edge.weight.rewardRef.rewardId;
        if(id && !containsReward(g.edge.weight.requirements, id) && contains(g.rewards, id)) {
            g.edge.weight.requirements.push({
                rewardId: id,
                quantity: (parseInt(g.edge.weight.rewardRef.quantity) || 1)//check if this can be removed if not set
            });
            resetRewardRef(g.edge.weight);
        }
    }

    g.addEdgeWeight = function() {
        var weight = g.edge.weight;
        if(!weight.rewardRef.edit) {
            g.edge.weights.push({
                requirements: weight.requirements,
                time: (parseInt(weight.time) || 1)//check if this can be removed if not set
            });
            resetEdgeWeight();
        }
    };

    g.addEdge = function() {
        var from = g.edge.from;
        var to = g.edge.to;
        if (from && to && contains(g.nodes, from) && contains(g.nodes, to)
                && !containsEdge(g.edges, from, to) && !g.edge.weight.edit) {
            var edge = {
                from: from,
                to: to,
                weights: g.edge.weights
            };
            g.edges.push(edge);
            networkEdges.add({
                id: g.edges.length-1,
                from: getNodeIndex(edge.from),
                to: getNodeIndex(edge.to),
                arrows: 'to',
                color: 'lightblue'
            });
            g.resetEdge();
        }
    };

    g.edgeRowSpan = function(weights) {
        var span = 0;
        for (var i = 0; i < g.len(weights); i++) {
            for (var j = 1; j < g.len(weights[i].requirements); j++) {
                span++;
            }
            span++;
        }
        return span;
    };

    g.removeReward = function(reward) {
        if(canRewardBeRemoved(reward.id)) {
            removeObj(g.rewards, reward);
        } else {
            log("Nonono");
        }
    };

    g.removeNode = function(node) {
        if(canNodeBeRemoved(node.id)) {
            networkNodes.remove({
                id: g.nodes.indexOf(node)
            });
            removeObj(g.nodes, node);
        } else {
            log("Nonono");
        }
    };

    g.removeEdge = function(edge) {
        networkEdges.remove({
            id: g.edges.indexOf(edge)
        });
        removeObj(g.edges, edge);
    };

    g.editReward = function(reward) {
        g.resetReward();
        toggleEdit(true, false, false);
        rewardBeingEdited = reward;
        g.reward.id = reward.id;
        g.reward.unique = reward.unique;
    };

    g.editNode = function(node) {
        g.resetNode();
        toggleEdit(false, true, false);
        nodeBeingEdited = node;
        g.node.id = node.id;
        g.node.revisitable = node.revisitable;
        g.node.rewards = angular.copy(node.rewards);
    };

    g.editEdge = function(edge) {
        g.resetEdge();
        toggleEdit(false, false, true);
        edgeBeingEdited = edge;
        g.edge.from = edge.from;
        g.edge.to = edge.to;
        g.edge.weights = angular.copy(edge.weights);
    };

    g.createGraph = function() {
        $http({
            method: 'POST',
            url: '/create/' + g.name + '/' + g.password,
            data: g.livesplit
        }).then(function successCallback(response) {
            $log.debug("SUCCeSS!");
        }, function errorCallback(response) {
            $log.debug("404 I guess...");
        });
    };

    resetColors = function() {
        for(var i = 0; i < networkNodes.length; i++) {
            networkNodes.update({
                id: i,
                color: 'lightblue'
            });
        }
        for(var i = 0; i < networkEdges.length; i++) {
            networkEdges.update({
                id: i,
                color: 'lightblue'
            });
        }
    }

    g.saveGraph = function() {
        resetColors();
        g.shortestPath = [];
        $http({
            method: 'POST',
            url: '/graph/' + g.name + '/' + g.password,
            data: {
                rewards: g.rewards,
                nodes: g.nodes,
                edges: g.edges,
                startId: g.startNode,
                endId: g.endNode
            }
        }).then(function successCallback(response) {
            var data = response.data;
            if(!data) {
                return;
            }
            g.shortestPath = data;
            for(var i = 0; i < data.length; i++) {
                if(i != 0) {
                    networkEdges.update({
                        id: getEdgeIndex(data[i-1], data[i]),
                        color: 'lime'
                    });
                }
                networkNodes.update({
                    id: getNodeIndex(data[i]),
                    color: 'lime'
                });
            }
        }, function errorCallback(response) {
            $log.debug("404 I guess...");
        });
    }

    loadGraph = function() {
        $http({
            method: 'GET',
            url: '/graph/' + g.name
        }).then(function successCallback(response) {
            var data = response.data;
            if(data.rewards) {
                g.rewards = data.rewards;
            }
            if(data.nodes) {
                g.nodes = data.nodes;
            }
            if(data.edges) {
                g.edges = data.edges;
            }
            if(data.startId) {
                g.startNode = data.startId;
            }
            if(data.endId) {
                 g.endNode = data.endId;
            }
            for(var i = 0; i < g.nodes.length; i++) {
                networkNodes.add({
                    id: i,
                    label: g.nodes[i].id,
                    color: 'lightblue'
                });
            }
            for(var i = 0; i < g.edges.length; i++) {
                networkEdges.add({
                    id: i,
                    from: getNodeIndex(g.edges[i].from),
                    to: getNodeIndex(g.edges[i].to),
                    arrows: 'to',
                    color: 'lightblue'
                });
            }
        }, function errorCallback(response) {
            $log.debug("404 I guess...");
        });
    }

    init = function() {
        var url = window.location.href;
        var parameter = 'g';
        var regex = new RegExp("[?&]" + parameter + "(=([^&#]*)|&|#|$)", "i");
        var results = regex.exec(url);
        if (!results || !results[2]) {
            log("Couldn't find key");
            return
        }
        g.name = decodeURIComponent(results[2].replace(/\+/g, " "));
        loadGraph()
    }

    g.resetReward();
    g.resetNode();
    g.resetEdge();
    init();
});